package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requests      = make(map[string]int)
	mutex         = &sync.Mutex{}
	maxRequests   = 10
	resetInterval = 1 * time.Minute
)

func BruteForceCheck(c *gin.Context) (bool, string) {
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	mutex.Lock()
	defer mutex.Unlock()

	count, exists := requests[ip]
	if !exists {
		requests[ip] = 1
		time.AfterFunc(resetInterval, func() {
			mutex.Lock()
			delete(requests, ip)
			mutex.Unlock()
		})
		return false, ""
	}

	if count >= maxRequests {
		return true, "Too many requests"
	}

	requests[ip]++
	return false, ""
}

func BruteForceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if blocked, msg := BruteForceCheck(c); blocked {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": msg})
			c.Abort()
			return
		}
		c.Next()
	}
} 