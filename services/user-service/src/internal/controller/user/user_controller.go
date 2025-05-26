package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
	"user-service/internal/models"
	"user-service/internal/user"
)

type CreateResponse struct {
	UserId string `json:"user_id"`
}

type CreateRequest struct {
	ChatId int64  `json:"chat_id"`
	Email  string `json:"email"`
	UUID   string `gorm:"id"`
}

type ShortUserResource struct {
	UUID   string `json:"uuid"`
	Email  string `json:"email"`
	ChatId int64  `json:"chat_id"`
}

type GetShortUserResponse struct {
	User ShortUserResource `json:"user"`
}

type GetUserResponse struct {
	User models.User `json:"user"`
}

type GetUserRequest struct {
	UUID string `json:"uuid"`
}

type SyncUsersRequest struct {
	UUIDs []string `json:"uuids"`
}

// Create godoc
// @Summary      Создать пользователя
// @Description  Создает нового пользователя с указанными данными и возвращает его идентификатор
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body   CreateRequest  true  "Данные пользователя"
// @Success      200   {object}  CreateResponse
// @Failure      400   {object}  object  "Неверный формат запроса"  example({"error": "invalid request"})
// @Failure      500   {object}  object  "Неверный формат запроса"  example({"error": "invalid request"})
// @Router       /user/create [post]
func Create(userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateRequest
		err := c.ShouldBind(&request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u := models.User{
			Email:     request.Email,
			UUID:      request.UUID,
			ChatId:    request.ChatId,
			CreatedAt: time.Time{},
		}
		id, err := userService.CreateUser(&u)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var response = CreateResponse{
			UserId: id,
		}

		c.JSON(200, gin.H{"result": response})
	}
}

// GetShortUser godoc
// @Summary      Получить укороченную информацию о пользователе
// @Description  Возвращает укороченную информацию о пользователе по его UUID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uuid  path  string  true  "UUID пользователя"
// @Success      200  {object}  GetShortUserResponse
// @Failure      400  {object}  object  "Неверный запрос"  example({"error": "invalid request"})
// @Failure      404  {object}  object  "Пользователь не найден" "{"error": "user not found"}" example(string)
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"  example({"error": "internal server error"})
// @Router       /user/{uuid}/short [get]
func GetShortUser(userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
			return
		}

		u, err := userService.GetByUUID(uuid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		response := GetShortUserResponse{
			User: ShortUserResource{
				UUID:   u.UUID,
				Email:  u.Email,
				ChatId: u.ChatId,
			},
		}

		c.JSON(http.StatusOK, gin.H{"result": response})
	}
}

// GetUser godoc
// @Summary      Получить пользователя
// @Description  Возвращает информацию о пользователе по его UUID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uuid  path  string  true  "UUID пользователя"
// @Success      200  {object}  GetShortUserResponse
// @Failure      400  {object}  object  "Неверный запрос"  example({"error": "invalid request"})
// @Failure      404  {object}  object  "Пользователь не найден" "{"error": "user not found"}" example(string)
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"  example({"error": "internal server error"})
// @Router       /user/{uuid} [get]
func GetUser(userService *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
			return
		}

		u, err := userService.GetByUUID(uuid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		response := GetUserResponse{
			User: *u,
		}

		c.JSON(http.StatusOK, gin.H{"result": response})
	}
}

// GetUserByChatId godoc
// @Summary      Получить пользователя по Chat ID
// @Description  Возвращает краткую информацию о пользователе по его Chat ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        chatId  path  string  true  "ChatID пользователя"
// @Success      200  {object}  ShortUserResource  "Успешный ответ"  example({"UUID": "550e8400-e29b-41d4-a716-446655440000", "Email": "user@example.com", "ChatId": 123456789})
// @Failure      400  {object}  object  "Неверный запрос"  example({"error": "invalid request"})
// @Failure      404  {object}  object  "Пользователь не найден"  example({"error": "user not found"})
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"  example({"error": "internal server error"})
// @Router       /user/by-chat/{chatId} [get]
func GetUserByChatId(userService *user.Service) gin.HandlerFunc {
	return func(context *gin.Context) {
		chatId, err := strconv.ParseInt(context.Param("chatId"), 10, 64)
		if err != nil {
			PP.Error(err.Error())
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if chatId == 0 {
			PP.Error(err.Error())
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
			return
		}
		u, err := userService.GetByChatId(chatId)
		if err != nil {
			PP.Error(err.Error())
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if u == nil {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		context.JSON(http.StatusOK, ShortUserResource{
			UUID:   u.UUID,
			Email:  u.Email,
			ChatId: u.ChatId,
		})
	}
}

// DeleteUserByChatId godoc
// @Summary      Удалить пользователя по Chat ID
// @Description  Удаляет пользователя из системы по его Chat ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        chatId  path  string  true  "ChatID пользователя"
// @Success      204  "Пользователь успешно удалён"
// @Failure      400  {object}  object  "Неверный запрос"  example({"error": "invalid request"})
// @Failure      404  {object}  object  "Пользователь не найден"  example({"error": "user not found"})
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"  example({"error": "internal server error"})
// @Router       /user/by-chat/{chatId} [delete]
func DeleteUserByChatId(userService *user.Service) gin.HandlerFunc {
	return func(context *gin.Context) {
		uuid := context.Param("uuid")

		if uuid == "" {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
			return
		}

		err := userService.DeleteByUUID(uuid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			PP.Error(err.Error())
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, gin.H{"result": fmt.Sprintf("deleted:%s", uuid)})
	}
}

func SyncUsers(userService *user.Service) gin.HandlerFunc {
	return func(context *gin.Context) {
		var request SyncUsersRequest
		err := context.ShouldBind(&request)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(request.UUIDs) <= 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
			return
		}
		synced, err := userService.SyncUsers(request.UUIDs)
		if err != nil {
			return
		}
		context.JSON(http.StatusOK, gin.H{"result": synced})
	}
}
