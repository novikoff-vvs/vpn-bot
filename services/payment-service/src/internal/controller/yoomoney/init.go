package yoomoney

import (
	"github.com/gin-gonic/gin"
	"github.com/novikoff-vvs/logger"
	"payment-service/internal/payment"
	"pkg/infrastructure/client/user"
	"pkg/infrastructure/http"
)

func RegisterRoutes(s *http.Server, client *user.Client, service *payment.Service, log logger.Interface, paymentSecret string) {
	registerWeb(s.GetWebGroup(), client)
	registerApi(s.GetApiGroup(), client, service, log, paymentSecret)
}

func registerWeb(r *gin.RouterGroup, client *user.Client) {
	yoomoneyEndpoints := r.Group("/yoomoney")
	{
		yoomoneyEndpoints.GET("/:uuid", paymentForm(client))
	}
}

func registerApi(r *gin.RouterGroup, client *user.Client, service *payment.Service, log logger.Interface, paymentSecret string) {
	yoomoneyEndpoints := r.Group("/yoomoney")
	{
		yoomoneyEndpoints.POST("/create", CreateYoomoneyLog(client, service, log, paymentSecret))
	}
}
