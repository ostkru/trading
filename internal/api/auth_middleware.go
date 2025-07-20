package api

import (
	"portaldata-api/internal/models"
	"portaldata-api/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	service *services.UserService
}

func NewAuthService(service *services.UserService) *AuthService {
	return &AuthService{service: service}
}

func BruteForceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if blocked, msg := BruteForceCheck(c); blocked {
			c.JSON(429, gin.H{"error": msg})
			c.Abort()
			return
		}
		c.Next()
	}
}

func (a *AuthService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("api_key")
		if apiKey == "" {
			apiKey = c.GetHeader("Authorization")
			if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
				apiKey = apiKey[7:]
			}
		}
		userID, err := a.service.GetUserIDByAPIKey(apiKey)
		if err != nil || userID == 0 {
			RegisterBruteForceFailure(c)
			c.JSON(401, models.APIResponse{
				OK:    false,
				Error: "Неверный API-ключ",
				Code:  401,
			})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
} 