package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"user-service/internal/models"
	"user-service/internal/repository/sqlite"
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

type GetUserResponse struct {
	User ShortUserResource `json:"user"`
}

type GetUserRequest struct {
	UUID string `json:"uuid"`
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
func Create(userRepo *sqlite.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request CreateRequest
		err := c.ShouldBind(&request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := models.User{
			Email:     request.Email,
			UUID:      request.UUID,
			ChatId:    request.ChatId,
			CreatedAt: time.Time{},
			IsActive:  true,
		}
		id, err := userRepo.Create(&user)
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

// GetUser godoc
// @Summary      Получить пользователя
// @Description  Возвращает информацию о пользователе по его UUID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uuid  path  string  true  "UUID пользователя"
// @Success      200  {object}  GetUserResponse
// @Failure      400  {object}  object  "Неверный запрос"  example({"error": "invalid request"})
// @Failure      404  {object}  object  "Пользователь не найден" "{"error": "user not found"}" example(string)
// @Failure      500  {object}  object  "Внутренняя ошибка сервера"  example({"error": "internal server error"})
// @Router       /user/{uuid} [get]
func GetUser(userRepo *sqlite.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		if uuid == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
			return
		}

		user, err := userRepo.GetByUUID(uuid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		response := GetUserResponse{
			User: ShortUserResource{
				UUID:   user.UUID,
				Email:  user.Email,
				ChatId: user.ChatId,
			},
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
func GetUserByChatId(userRepo *sqlite.UserRepository) gin.HandlerFunc {
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
		user, err := userRepo.GetByChatId(chatId)
		if err != nil {
			PP.Error(err.Error())
			context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user == nil {
			PP.Error(err.Error())
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		context.JSON(http.StatusOK, ShortUserResource{
			UUID:   user.UUID,
			Email:  user.Email,
			ChatId: user.ChatId,
		})
	}
}
