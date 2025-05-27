package yoomoney

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"payment-service/internal/model"
	"payment-service/internal/payment"
	"payment-service/internal/singleton"
	"pkg/exceptions"
	"pkg/infrastructure/client/subscription"
	"pkg/infrastructure/client/user"
	response2 "pkg/infrastructure/client/user/response"
	singleton2 "pkg/singleton"
	"strings"
	"time"
)

type LogRequest struct {
	NotificationType string    `form:"notification_type" json:"notification_type"`
	OperationId      string    `form:"operation_id" json:"operation_id"`
	Amount           float32   `form:"amount" json:"amount"`
	WithdrawAmount   float32   `form:"withdraw_amount" json:"withdraw_amount"`
	Datetime         time.Time `form:"datetime" json:"datetime"`
	Label            string    `form:"label" json:"label"`
}

func CreateYoomoneyLog(client *user.Client, service *payment.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer c.JSON(http.StatusOK, gin.H{})
		var request LogRequest

		if err := c.ShouldBind(&request); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err)
			c.Abort()
			return
		}

		if request.Label == "" {
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		_, err := client.GetUserByUUID(user.GetUserByUUIDRequest{UUID: request.Label})
		if err != nil {
			//TODO: добавить обработку ненахождения юзера
			log.Println(err)
			return
		}

		decrypted, err := singleton.CryptoService().Decrypt(request.Label)
		if err != nil {
			//TODO: добавить обработку
			log.Println(err)
			return
		}

		decryptedString := strings.Split(string(decrypted), ":")
		if len(decryptedString) != 2 {
			//TODO: добавить обработку проблем с токеном
			return
		}
		uuid := decryptedString[0]

		var p = model.Payment{
			PaymentAmount:  float64(request.Amount),
			WithdrawAmount: float64(request.WithdrawAmount),
			OperationId:    request.OperationId,
			CreatedAt:      request.Datetime,
			Token:          uuid,
		}

		err = service.LogPayment(p)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = singleton2.SubscriptionClient().GetSubscriptionByUUID(subscription.GetSubscriptionByUUIDRequest{UUID: uuid})
		if err != nil {
			return
		}
	}
}

func paymentForm(client *user.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		response, err := client.GetUserByUUID(user.GetUserByUUIDRequest{UUID: uuid})
		if errors.Is(err, exceptions.ErrModelNotFound) {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var result response2.GetUserByUUIDResponse

		err = json.Unmarshal(response.Bytes(), &result)
		if err != nil {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		cryptoService := singleton.CryptoService()
		encStr := result.Result.User.UUID + ":" + time.Now().String()
		encrypted, err := cryptoService.Encrypt(encStr)
		if err != nil {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.HTML(http.StatusOK, "payload.html", gin.H{
			"label":   encrypted,
			"tg_name": result.Result.User.Email,
			"amount":  result.Result.User.Subscription.Plan.Price,
		})
	}
}
