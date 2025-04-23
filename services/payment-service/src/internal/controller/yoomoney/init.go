package yoomoney

import (
	"github.com/gin-gonic/gin"
	"pkg/infrastructure/http"
)

func RegisterRoutes(s *http.Server) {
	registerWeb(s.GetWebGroup())
}

func registerWeb(r *gin.RouterGroup) {
	yoomoneyEndpoints := r.Group("/yoomoney")
	{

		yoomoneyEndpoints.GET("/:uuid", paymentForm)

		//yoomoneyEndpoints.POST("/create", createYoomoneyLog)
		//yoomoneySecureEndpoints := yoomoneyEndpoints.Group("/secure")
		////yoomoneySecureEndpoints.Use(AuthRequired())
		//{
		//	yoomoneySecureEndpoints.GET("/", indexYoomoney)
		//}
	}
}
