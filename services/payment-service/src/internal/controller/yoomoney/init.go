package yoomoney

import (
	"github.com/gin-gonic/gin"
	"payment-service/internal/payment"
	"pkg/infrastructure/client/user"
	"pkg/infrastructure/http"
)

func RegisterRoutes(s *http.Server, client *user.Client, service *payment.Service) {
	registerWeb(s.GetWebGroup(), client)
	registerApi(s.GetApiGroup(), client, service)
}

func registerWeb(r *gin.RouterGroup, client *user.Client) {
	yoomoneyEndpoints := r.Group("/yoomoney")
	{
		yoomoneyEndpoints.GET("/:uuid", paymentForm(client))

		//yoomoneyEndpoints.POST("/create", CreateYoomoneyLog)
		//yoomoneySecureEndpoints := yoomoneyEndpoints.Group("/secure")
		////yoomoneySecureEndpoints.Use(AuthRequired())
		//{
		//	yoomoneySecureEndpoints.GET("/", indexYoomoney)
		//}
	}
}

func registerApi(r *gin.RouterGroup, client *user.Client, service *payment.Service) {
	yoomoneyEndpoints := r.Group("/yoomoney")
	{
		yoomoneyEndpoints.POST("/create", CreateYoomoneyLog(client, service))
	}
}
