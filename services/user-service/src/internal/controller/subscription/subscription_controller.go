package subscription

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"pkg/events"
	"pkg/singleton"
	"user-service/internal/subscription"
)

type CreateRequest struct {
	UUID string `gorm:"id"`
}

type RefreshRequest struct {
	UserUUID      string `json:"user_uuid"`
	AmountPeriods int    `json:"amount_periods"`
}

func refresh(subscriptionService *subscription.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RefreshRequest
		err := c.BindJSON(&request)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		s, err := subscriptionService.Refresh(subscription.RefreshDTO{
			UserUUID:      request.UserUUID,
			AmountPeriods: request.AmountPeriods,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		marshal, err := json.Marshal(events.SubscriptionRefreshed{
			UserUUID: s.UserUUID,
			ChatId:   s.User.ChatId,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		err = singleton.NatsPublisher().Publish("events.subscription.refreshed", marshal)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": s})
	}
}

func getSubscriptionByUser(subscriptionService *subscription.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("id")
		getSubscription, err := subscriptionService.GetActiveSubscriptionByUser(uuid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": getSubscription})
	}
}
