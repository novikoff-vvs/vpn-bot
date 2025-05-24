package yoomoney

import (
	"github.com/gin-gonic/gin"
	"pkg/infrastructure/client/user"
	"pkg/infrastructure/http"
)

func RegisterRoutes(s *http.Server, client *user.Client) {
	registerWeb(s.GetWebGroup(), client)
}

func registerWeb(r *gin.RouterGroup, client *user.Client) {
	yoomoneyEndpoints := r.Group("/yoomoney")
	{
		yoomoneyEndpoints.GET("/:uuid", paymentForm(client))

		yoomoneyEndpoints.GET("/:uuid", paymentForm)

		//yoomoneyEndpoints.POST("/create", createYoomoneyLog)
		//yoomoneySecureEndpoints := yoomoneyEndpoints.Group("/secure")
		////yoomoneySecureEndpoints.Use(AuthRequired())
		//{
		//	yoomoneySecureEndpoints.GET("/", indexYoomoney)
		//}
	}
}
