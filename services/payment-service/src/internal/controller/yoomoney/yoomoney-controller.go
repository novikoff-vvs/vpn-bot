package yoomoney

import (
	"github.com/gin-gonic/gin"
	"payment-service/internal/infrastructure/http"
)

func RegisterRoutes(s *http.Server) {
	registerApi(s.GetApiGroup())
	registerWeb(s.GetWebGroup())
}

func registerApi(group *gin.RouterGroup) {

}

func registerWeb(group *gin.RouterGroup) {
	yoomoneyEndpoints := group.Group("/yoomoney")
	{

		yoomoneyEndpoints.GET("/", paymentForm)

		yoomoneyEndpoints.POST("/create", createYoomoneyLog)
		yoomoneySecureEndpoints := yoomoneyEndpoints.Group("/secure")
		yoomoneySecureEndpoints.Use(AuthRequired())
		{
			yoomoneySecureEndpoints.GET("/", indexYoomoney)
		}
	}
}
