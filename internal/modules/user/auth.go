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
		if authHeader != "" && len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = authHeader[7:]
		} else {
			apiKey = c.GetHeader("X-API-KEY")
		}
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется API ключ (используйте Authorization: Bearer <токен> или X-API-KEY)"})
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
