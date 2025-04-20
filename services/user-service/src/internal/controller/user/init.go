package user

import (
	"github.com/gin-gonic/gin"
	"github.com/novikoff-vvs/logger"
	"pkg/infrastructure/http"
	"user-service/internal/repository/sqlite"
)

var PP logger.Interface

func RegisterRoutes(s *http.Server, userRepo *sqlite.UserRepository, p logger.Interface) {
	PP = p
	registerApi(s.GetApiGroup(), userRepo)
}

func registerApi(r *gin.RouterGroup, userRepo *sqlite.UserRepository) {
	group := r.Group("/user")
	{
		group.POST("/create", Create(userRepo))
		groupByUUID := group.Group(":uuid")
		{
			groupByUUID.GET("", GetUser(userRepo))
		}

		groupByChatId := group.Group("by-chat")
		{
			groupByChatId.GET(":chatId", GetUserByChatId(userRepo))
		}
	}
}
