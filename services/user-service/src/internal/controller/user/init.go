package user

import (
	"github.com/gin-gonic/gin"
	"github.com/novikoff-vvs/logger"
	"pkg/infrastructure/http"
	"user-service/internal/user"
)

var PP logger.Interface //todo нужно нахуй это выпилить и использовать логгер как-то по-другому или хотя бы переименовать

func RegisterRoutes(s *http.Server, userService *user.Service, p logger.Interface) {
	PP = p
	registerApi(s.GetApiGroup(), userService)
}

func registerApi(r *gin.RouterGroup, userService *user.Service) {
	group := r.Group("/user")
	{
		group.GET("/all", All(userService))
		group.POST("/create", Create(userService))
		group.POST("/sync-users", SyncUsers(userService))

		groupByUUID := group.Group(":uuid")
		{
			groupByUUID.GET("/short", GetShortUser(userService))
			groupByUUID.GET("", GetUser(userService))
			groupByUUID.DELETE("", DeleteUserByChatId(userService))
		}

		groupByChatId := group.Group("by-chat/:chatId")
		{
			groupByChatId.GET("", GetUserByChatId(userService))
		}
	}
}
