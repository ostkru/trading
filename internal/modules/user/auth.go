package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	userService *Service
}

func NewAuthService(userService *Service) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var apiKey string

		// Проверяем заголовок Authorization (Bearer token)
		if authHeader != "" && len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = authHeader[7:]
		} else {
			// Проверяем заголовок X-API-KEY
			apiKey = c.GetHeader("X-API-KEY")
		}

		// Если API ключ не найден в заголовках, проверяем GET параметр api_key
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется API ключ (используйте Authorization: Bearer <токен>, X-API-KEY или api_key в GET параметрах)"})
			c.Abort()
			return
		}

		userID, err := s.userService.GetUserIDByAPIKey(apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный API ключ"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
