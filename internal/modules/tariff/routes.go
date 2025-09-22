package tariff

import "github.com/gin-gonic/gin"

// RegisterRoutes регистрирует маршруты для тарифов
func RegisterRoutes(router *gin.RouterGroup, handler *Handler) {
	tariffs := router.Group("/tariffs")
	{
		// Публичные маршруты (только просмотр активных тарифов)
		tariffs.GET("", handler.ListTariffs)        // GET /tariffs
		tariffs.GET("/:id", handler.GetTariff)      // GET /tariffs/{id}
		tariffs.GET("/my", handler.GetMyTariffInfo) // GET /tariffs/my

		// Административные маршруты (требуют специальных прав)
		tariffs.POST("", handler.CreateTariff)                   // POST /tariffs
		tariffs.PUT("/:id", handler.UpdateTariff)                // PUT /tariffs/{id}
		tariffs.DELETE("/:id", handler.DeleteTariff)             // DELETE /tariffs/{id}
		tariffs.GET("/user/:user_id", handler.GetUserTariffInfo) // GET /tariffs/user/{user_id}
		tariffs.PUT("/user/change", handler.ChangeUserTariff)    // PUT /tariffs/user/change
		tariffs.GET("/stats", handler.GetTariffUsageStats)       // GET /tariffs/stats
	}
}
