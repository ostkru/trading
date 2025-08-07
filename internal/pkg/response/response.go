package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse представляет успешный ответ API
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Success возвращает успешный ответ с данными
func Success(c *gin.Context, statusCode int, data interface{}, message string) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	c.JSON(statusCode, response)
}

// SuccessWithData возвращает успешный ответ только с данными
func SuccessWithData(c *gin.Context, statusCode int, data interface{}) {
	Success(c, statusCode, data, "")
}

// SuccessWithMessage возвращает успешный ответ только с сообщением
func SuccessWithMessage(c *gin.Context, statusCode int, message string) {
	Success(c, statusCode, nil, message)
}

// Error возвращает ответ с ошибкой
func Error(c *gin.Context, statusCode int, errorMessage string) {
	response := ErrorResponse{
		Success: false,
		Error:   errorMessage,
	}
	c.JSON(statusCode, response)
}

// BadRequest возвращает ошибку 400
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized возвращает ошибку 401
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// Forbidden возвращает ошибку 403
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}

// NotFound возвращает ошибку 404
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalServerError возвращает ошибку 500
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
} 