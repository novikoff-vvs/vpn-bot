package yoomoney

import (
	"errors"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"strconv"
	"time"
)

type yoomoneyLogRequest struct {
	NotificationType string    `form:"notification_type"`
	OperationId      string    `form:"operation_id"`
	Amount           float32   `form:"amount"`
	WithdrawAmount   float32   `form:"withdraw_amount"`
	Datetime         time.Time `form:"datetime"`
	Label            string    `form:"label"`
}

func createYoomoneyLog(c *gin.Context) {
	var request yoomoneyLogRequest
	var settings BotSettings

	DB.Model(&BotSettings{}).First(&settings)

	if err := c.ShouldBind(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		return
	}

	if request.Label == "" {
		c.JSON(http.StatusOK, gin.H{})
	}
	var payment UserPayment
	DB.Model(&UserPayment{}).Preload(clause.Associations).Where("uuid = ?", request.Label).Where("uuid != ''").First(&payment)
	var newPayment UserPayment
	err := DB.Transaction(func(tx *gorm.DB) error {
		var err error
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
		if payment.ID == 0 {
			if settings.AdminUserId != 0 {
				_, err = Bot.Send(tgbotapi.NewMessage(settings.AdminUserId, "Не удалось найти payment.ID. "+request.OperationId))
				if err != nil {
					log.Println(err)
				}
			}
			return errors.New("Не найден payment.UUID " + request.Label)
		}

		payment.IsPaid = true
		if err = tx.Save(&payment).Error; err != nil {
			// Если произошла ошибка, откатываем транзакцию
			return err
		}
		//TODO вынести сеттинги в общий прелоад
		if settings.AdminUserId != 0 {
			_, err = Bot.Send(tgbotapi.NewMessage(settings.AdminUserId, "Пользователь: "+strconv.Itoa(int(payment.VpnUser.ID))+" произвел оплату"))
			if err != nil {
				log.Println(err)
			}
		}

		yoomoneyLog := YoomoneyLog{
			NotificationType: request.NotificationType,
			OperationId:      request.OperationId,
			Amount:           request.Amount,
			WithdrawAmount:   request.WithdrawAmount,
			Datetime:         request.Datetime,
			Label:            request.Label,
		}
		if err = tx.Save(&yoomoneyLog).Error; err != nil {
			// Если произошла ошибка, откатываем транзакцию
			return err
		}

		uid, err := uuid.NewRandom()
		if err != nil {
			return err
		}

		newPayment = UserPayment{
			Model:         gorm.Model{},
			PaymentAmount: payment.PaymentAmount,
			VpnUserID:     payment.VpnUserID,
			IsPaid:        false,
			InvoiceDate:   getNextPaymentDate(payment.InvoiceDate),
			UUID:          uid.String(),
		}
		if err = tx.Save(&newPayment).Error; err != nil {
			// Если произошла ошибка, откатываем транзакцию
			return err
		}

		return nil
	})

	if err != nil {
		if settings.AdminUserId != 0 {
			_, err = Bot.Send(tgbotapi.NewMessage(settings.AdminUserId, "Автоматическая оплата не прошла! UUID: "+request.Label+" : "+request.OperationId))
			if err != nil {
				log.Println(err)
			}
		}
		if payment.ID != 0 {
			err = BotSendMessage(
				tgbotapi.NewMessage(
					payment.VpnUser.ChatId,
					"Платеж получен, но не обработан автоматически. Обратитесь к @code_style и сообщите UUID: "+request.Label,
				))
			if err != nil {
				log.Println(err)
			}
		}
	}

	if newPayment.ID != 0 {
		err = BotSendMessage(tgbotapi.NewMessage(payment.VpnUser.ChatId, "Дата следующей оплаты: "+newPayment.InvoiceDate.Format("02-01-2006")))
		if err != nil {
			log.Println(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func paymentForm(c *gin.Context) {
	var payment UserPayment
	uuid := c.DefaultQuery("uuid", "null")
	DB.Model(&UserPayment{}).Preload(clause.Associations).Where("uuid = ?", uuid).First(&payment)

	if payment.IsPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment is paid"})
	}

	c.HTML(http.StatusOK, "payload.html", gin.H{
		"label":   payment.UUID, // Данные для передачи в шаблон
		"tg_name": payment.VpnUser.Username,
		"amount":  FindAmountWithCommission(payment.PaymentAmount),
	})
}

func indexYoomoney(c *gin.Context) {
	var yoomoneys []YoomoneyLog
	DB.Model(&YoomoneyLog{}).Preload(clause.Associations).Find(&yoomoneys)
	c.JSON(http.StatusOK, yoomoneys)
}
