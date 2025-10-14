package user

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	service *Service
}

func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// RegisterRoutes регистрирует маршруты пользователей
func RegisterRoutes(router *gin.RouterGroup, handlers *Handlers) {
	// Публичные маршруты (без авторизации)
	router.POST("/user/registration", handlers.RegisterUser)
	router.POST("/authentication_token", handlers.AuthenticateUser)
	router.POST("/user/create_password_reset", handlers.CreatePasswordReset)
	router.POST("/user/password_reset/:token", handlers.ResetPassword)
	router.POST("/user/verify", handlers.VerifyUser)
	
	// Защищенные маршруты (с авторизацией)
	router.GET("/user/data", handlers.GetUserData)
	router.PUT("/user", handlers.UpdateUser)
	router.DELETE("/user", handlers.DeleteUser)
}

// RegisterUserRequest запрос на регистрацию
type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterUser регистрирует нового пользователя
func (h *Handlers) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Неверные данные: " + err.Error(),
		})
		return
	}

	// Генерируем API токен
	apiToken := generateAPIToken()

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка хеширования пароля",
		})
		return
	}

	// Создаем пользователя
	createReq := &CreateUserRequest{
		Username: req.Name,
		Email:    req.Email,
		APIToken: apiToken,
		TariffID: 1, // Базовый тариф
	}

	user, err := h.service.CreateUser(createReq)
	if err != nil {
		if strings.Contains(err.Error(), "уже существует") {
			c.JSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Ошибка создания пользователя: " + err.Error(),
			})
		}
		return
	}

	// Сохраняем хешированный пароль в отдельной таблице
	_, err = h.service.db.Exec("INSERT INTO user_passwords (user_id, password_hash) VALUES (?, ?)", user.ID, string(hashedPassword))
	if err != nil {
		// Если не удалось сохранить пароль, удаляем пользователя
		h.service.DeleteUser(user.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка сохранения пароля",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Пользователь успешно зарегистрирован",
		"data": gin.H{
			"user_id": user.ID,
			"api_token": user.APIToken,
		},
	})
}

// AuthenticateRequest запрос на аутентификацию
type AuthenticateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthenticateUser аутентифицирует пользователя
func (h *Handlers) AuthenticateUser(c *gin.Context) {
	var req AuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Неверные данные: " + err.Error(),
		})
		return
	}

	// Получаем пользователя по email
	var userID int64
	var passwordHash string
	err := h.service.db.QueryRow("SELECT u.id, up.password_hash FROM users u JOIN user_passwords up ON u.id = up.user_id WHERE u.email = ? AND u.is_active = TRUE", req.Email).Scan(&userID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Неверный email или пароль",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Ошибка аутентификации",
			})
		}
		return
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Неверный email или пароль",
		})
		return
	}

	// Получаем пользователя
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка получения данных пользователя",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Успешная аутентификация",
		"data": gin.H{
			"token": user.APIToken,
			"user": user,
		},
	})
}

// CreatePasswordResetRequest запрос на создание токена сброса пароля
type CreatePasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// CreatePasswordReset создает токен для сброса пароля
func (h *Handlers) CreatePasswordReset(c *gin.Context) {
	var req CreatePasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Неверные данные: " + err.Error(),
		})
		return
	}

	// Проверяем существование пользователя
	var userID int64
	err := h.service.db.QueryRow("SELECT id FROM users WHERE email = ? AND is_active = TRUE", req.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "Пользователь с таким email не найден",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Ошибка поиска пользователя",
			})
		}
		return
	}

	// Генерируем токен сброса пароля
	token := generateAPIToken()
	expiresAt := time.Now().Add(24 * time.Hour) // Токен действителен 24 часа

	// Сохраняем токен
	_, err = h.service.db.Exec("INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка создания токена сброса пароля",
		})
		return
	}

	// TODO: Отправить email с токеном
	// Пока просто возвращаем токен для тестирования
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Токен сброса пароля создан",
		"data": gin.H{
			"token": token,
			"expires_at": expiresAt,
		},
	})
}

// ResetPasswordRequest запрос на сброс пароля
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// ResetPassword сбрасывает пароль по токену
func (h *Handlers) ResetPassword(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Токен не указан",
		})
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Неверные данные: " + err.Error(),
		})
		return
	}

	// Проверяем токен
	var userID int64
	var expiresAt time.Time
	err := h.service.db.QueryRow("SELECT user_id, expires_at FROM password_reset_tokens WHERE token = ? AND expires_at > NOW()", token).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Недействительный или истекший токен",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Ошибка проверки токена",
			})
		}
		return
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка хеширования пароля",
		})
		return
	}

	// Обновляем пароль
	_, err = h.service.db.Exec("UPDATE user_passwords SET password_hash = ? WHERE user_id = ?", string(hashedPassword), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка обновления пароля",
		})
		return
	}

	// Удаляем использованный токен
	_, err = h.service.db.Exec("DELETE FROM password_reset_tokens WHERE token = ?", token)
	if err != nil {
		// Логируем ошибку, но не прерываем процесс
		fmt.Printf("Ошибка удаления токена: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Пароль успешно изменен",
	})
}

// VerifyRequest запрос на верификацию
type VerifyRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyUser верифицирует пользователя
func (h *Handlers) VerifyUser(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Неверные данные: " + err.Error(),
		})
		return
	}

	// Проверяем токен верификации
	var userID int64
	err := h.service.db.QueryRow("SELECT user_id FROM verification_tokens WHERE token = ? AND expires_at > NOW()", req.Token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Недействительный или истекший токен",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Ошибка проверки токена",
			})
		}
		return
	}

	// Обновляем статус верификации пользователя
	_, err = h.service.db.Exec("UPDATE users SET is_verified = TRUE WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка обновления статуса верификации",
		})
		return
	}

	// Удаляем использованный токен
	_, err = h.service.db.Exec("DELETE FROM verification_tokens WHERE token = ?", req.Token)
	if err != nil {
		fmt.Printf("Ошибка удаления токена верификации: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Пользователь успешно верифицирован",
	})
}

// GetUserData получает данные текущего пользователя
func (h *Handlers) GetUserData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Пользователь не авторизован",
		})
		return
	}

	user, err := h.service.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка получения данных пользователя",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Данные пользователя получены",
		"item": user,
	})
}

// UpdateUser обновляет данные пользователя
func (h *Handlers) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Пользователь не авторизован",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Неверные данные: " + err.Error(),
		})
		return
	}

	user, err := h.service.UpdateUser(userID.(int64), &req)
	if err != nil {
		if strings.Contains(err.Error(), "уже существует") {
			c.JSON(http.StatusConflict, gin.H{
				"code":    409,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Ошибка обновления пользователя: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Пользователь успешно обновлен",
		"data": user,
	})
}

// DeleteUser удаляет пользователя
func (h *Handlers) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "Пользователь не авторизован",
		})
		return
	}

	err := h.service.DeleteUser(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Ошибка удаления пользователя: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "Пользователь успешно удален",
	})
}

// generateAPIToken генерирует случайный API токен
func generateAPIToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}
