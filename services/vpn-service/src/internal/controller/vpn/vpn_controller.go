package vpn

import (
	"net/http"
	"pkg/models"
	"strconv"
	"time"
	"vpn-service/internal/service/vpn"
	"vpn-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	ChatId int64  `json:"chat_id"`
	Email  string `json:"email"`
	UUID   string `json:"uuid"`
}

type CreateResponse struct {
	Message string `json:"message"`
}

type ShortUserResource struct {
	UUID   string `json:"uuid"`
	Email  string `json:"email"`
	ChatId int64  `json:"chat_id"`
}

// RegisterUser godoc
// @Summary      Зарегистрировать пользователя VPN
// @Description  Создает нового пользователя VPN и добавляет его в прокси-сервер
// @Tags         vpn
// @Accept       json
// @Produce      json
// @Param        request  body     CreateRequest  true  "Данные нового пользователя"
// @Success      200      {object} CreateResponse "Пользователь успешно зарегистрирован"
// @Failure      400      {object} object         "Неверный формат запроса"  example({"error": "invalid request"})
// @Failure      500      {object} object         "Ошибка сервера"            example({"error": "internal server error"})
// @Router       /vpn/register [post]
func RegisterUser(service vpn.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &models.VpnUser{
			ChatId: req.ChatId,
			Email:  req.Email,
			UUID:   req.UUID,
		}

		err := service.UserRegisterByChatId(user, "Registered via API at "+time.Now().Format(time.RFC3339))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, CreateResponse{Message: "user registered"})
	}
}

// GetUserByChatId godoc
// @Summary      Получить пользователя VPN по Chat ID
// @Description  Возвращает информацию о VPN-пользователе по его Telegram Chat ID
// @Tags         vpn
// @Accept       json
// @Produce      json
// @Param        chatId  path     string             true  "Chat ID пользователя"
// @Success      200     {object} ShortUserResource  "Информация о пользователе"
// @Failure      400     {object} object             "Неверный формат Chat ID" example({"error": "invalid chat id"})
// @Failure      404     {object} object             "Пользователь не найден" example({"error": "user not found"})
// @Router       /vpn/by-chat/{chatId} [get]
func GetUserByChatId(service vpn.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatIdStr := c.Param("chatId")
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
			return
		}

		user, err := service.UserGetByChatId(chatId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		res := ShortUserResource{
			UUID:   user.UUID,
			Email:  user.Email,
			ChatId: user.ChatId,
		}
		c.JSON(http.StatusOK, res)
	}
}

// ResetTraffic godoc
// @Summary      Сбросить трафик пользователя
// @Description  Обнуляет трафик VPN-пользователя по его Chat ID
// @Tags         vpn
// @Accept       json
// @Produce      json
// @Param        chatId  path     string  true  "Chat ID пользователя"
// @Success      200     {object} object  "Трафик успешно сброшен" example({"message": "traffic reset successfully"})
// @Failure      400     {object} object  "Неверный формат Chat ID" example({"error": "invalid chat id"})
// @Failure      500     {object} object  "Ошибка сервера" example({"error": "internal server error"})
// @Router       /vpn/reset-traffic/{chatId} [post]
func ResetTraffic(service vpn.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatIdStr := c.Param("chatId")
		chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
			return
		}

		if err := service.ResetClientTraffic(chatId); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "traffic reset successfully"})
	}
}

// UserExists godoc
// @Summary      Проверить существование пользователя
// @Description  Проверяет, существует ли пользователь с данным Chat ID
// @Tags         vpn
// @Produce      json
// @Param        chatId  path   int  true  "Chat ID"
// @Success      200   {object}  gin.H  "Пример: {\"exists\": true}"
// @Failure      400   {object}  gin.H
// @Router       /vpn/exists/{chatId} [get]
func UserExists(service vpn.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatId, err := strconv.ParseInt(c.Param("chatId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
			return
		}
		exists := service.UserExistsByChatId(chatId)
		c.JSON(http.StatusOK, gin.H{"exists": exists})
	}
}

// GetUserByEmail godoc
// @Summary      Получить пользователя по Email
// @Description  Возвращает пользователя по email-адресу
// @Tags         vpn
// @Produce      json
// @Param        email  path   string  true  "Email пользователя"
// @Success      200   {object}  ShortUserResource
// @Failure      404   {object}  gin.H
// @Router       /vpn/exists/{chatId} [get]
func GetUserByEmail(service vpn.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")
		user, err := service.UserGetByEmail(email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, ShortUserResource{
			UUID:   user.UUID,
			Email:  user.Email,
			ChatId: user.ChatId,
		})
	}
}

type SubscriptionLinkResource struct {
	SubscriptionLink string `json:"subscription_link"`
}

// GetSubcLinkByChatId godoc
// @Summary      Получить ссылку на подписку
// @Description  Возвращает ссылку на подписочную систему по пользователю
// @Tags         vpn
// @Produce      json
// @Param        chatId  path   string  true  "chatId пользователя"
// @Success      200   {object}  SubscriptionLinkResource
// @Failure      404   {object}  gin.H
// @Router       /vpn/subscription-link/{chatId}  [get]
func GetSubcLinkByChatId(service vpn.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatId, err := strconv.ParseInt(c.Param("chatId"), 10, 64)
		user, err := service.UserGetByChatId(chatId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		subscLink := utils.BuildVlessLink(user.UUID)
		c.JSON(http.StatusOK, SubscriptionLinkResource{
			SubscriptionLink: subscLink,
		})
	}
}
