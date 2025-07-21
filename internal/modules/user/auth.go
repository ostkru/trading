package user

import (
	"net/http"

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
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			c.Abort()
			return
		}

		userID, err := s.userService.GetUserIDByAPIKey(apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
} 