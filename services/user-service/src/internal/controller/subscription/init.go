package subscription

import (
	"github.com/gin-gonic/gin"
	"pkg/infrastructure/http"
	"user-service/internal/subscription"
)

func RegisterRoutes(s *http.Server, subscriptionService *subscription.Service) {
	registerApi(s.GetApiGroup(), subscriptionService)
}

func registerApi(r *gin.RouterGroup, subscriptionService *subscription.Service) {
	subscriptionEndpoints := r.Group("/subscription")
	{
		subscriptionEndpoints.POST("/refresh", refresh(subscriptionService))
		byUser := subscriptionEndpoints.Group("/by-user")
		{
			byUser.GET("/:id", getSubscriptionByUser(subscriptionService))
		}
	}
}
