package yoomoney

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/novikoff-vvs/logger"
	"net/http"
	"net/url"
	"payment-service/internal/model"
	"payment-service/internal/payment"
	"payment-service/internal/singleton"
	"pkg/exceptions"
	"pkg/infrastructure/client/user"
	response2 "pkg/infrastructure/client/user/response"
	"strconv"
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

func CreateYoomoneyLog(client *user.Client, service *payment.Service, log logger.Interface, paymentSecret string) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer c.JSON(http.StatusOK, gin.H{})

		var request LogRequest

		if err := c.Request.ParseForm(); err != nil {
			log.Error(fmt.Sprintf("Failed to parse form: %v", err))
			return
		}

		log.Info(fmt.Sprintf("%+v\n", c.Request.PostForm))

		for key, values := range c.Request.PostForm {
			for _, value := range values {
				log.Info(fmt.Sprintf("PostForm param: %s = %s", key, value))
			}
		}

		form := c.Request.PostForm
		form["notification_secret"] = []string{paymentSecret}

		keys := []string{
			"notification_type",
			"operation_id",
			"amount",
			"currency",
			"datetime",
			"sender",
			"codepro",
			"notification_secret",
			"label",
		}

		var params []string
		for _, key := range keys {
			if key == "sha1_hash" {
				continue
			}
			values := form[key]
			if len(values) > 0 {
				val := values[0]
				params = append(params, val)
			}
		}

		result := strings.Join(params, "&")
		log.Info(result)

		hash := sha1.New()
		hash.Write([]byte(result))
		sha1Hash := hex.EncodeToString(hash.Sum(nil))

		log.Info(fmt.Sprintf("Calculated SHA1 hash: %s", sha1Hash))
		log.Info(fmt.Sprintf("Requested SHA1 hash: %s", c.Request.FormValue("sha1_hash")))

		if c.Request.FormValue("sha1_hash") != sha1Hash {
			log.Error("SHA1 hash mismatch – possible spoofed request or secret mismatch.")
			//return
		}

		log.Info("SHA1 hash matched — notification is verified.")

		request.Label = c.PostForm("label")
		request.OperationId = c.PostForm("operation_id")
		request.NotificationType = c.PostForm("notification_type")

		datetimeStr := c.PostForm("datetime")
		parsedTime, err := time.Parse(time.RFC3339, datetimeStr)
		if err != nil {
			log.Error(fmt.Sprintf("Invalid datetime format: %s, error: %v", datetimeStr, err))
			return
		}
		request.Datetime = parsedTime

		amountStr := c.PostForm("amount")
		amount, err := strconv.ParseFloat(amountStr, 32)
		if err != nil {
			log.Error(fmt.Sprintf("Invalid amount format: %s, error: %v", amountStr, err))
			return
		}
		request.Amount = float32(amount)

		withdrawStr := c.PostForm("withdraw_amount")
		if withdrawStr == "" {
			log.Info("withdraw_amount is empty, defaulting to 0")
			request.WithdrawAmount = 0
		} else {
			withdraw, err := strconv.ParseFloat(withdrawStr, 32)
			if err != nil {
				log.Error(fmt.Sprintf("Invalid withdraw_amount format: %s, error: %v", withdrawStr, err))
				return
			}
			request.WithdrawAmount = float32(withdraw)
		}

		if request.Label == "" {
			log.Error("Empty label received")
			return
		}

		uuid := request.Label

		_, err = client.GetUserByUUID(user.GetUserByUUIDRequest{UUID: request.Label})
		if err != nil {
			log.Error(fmt.Sprintf("User not found with UUID: %s, error: %v", uuid, err))
			return
		}

		p := model.Payment{
			PaymentAmount:  float64(request.Amount),
			WithdrawAmount: float64(request.WithdrawAmount),
			OperationId:    request.OperationId,
			CreatedAt:      request.Datetime,
			Token:          uuid,
		}

		err = service.LogPayment(p)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to log payment: %v", err))
			return
		}
		log.Info(fmt.Sprintf("Payment logged: %+v", p))
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
		encStr := result.Result.User.UUID + "::" + strconv.FormatInt(time.Now().Unix(), 10)
		encrypted, err := cryptoService.Encrypt(encStr)
		if err != nil {
			err = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		url.QueryEscape(encrypted)

		c.HTML(http.StatusOK, "payload.html", gin.H{
			"label":   result.Result.User.UUID,
			"tg_name": result.Result.User.Email,
			"amount":  result.Result.User.Subscription.Plan.Price,
		})
	}
}
