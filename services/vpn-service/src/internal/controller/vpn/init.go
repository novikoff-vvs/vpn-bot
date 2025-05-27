package vpn

import (
	"github.com/gin-gonic/gin"
	"pkg/infrastructure/http"
	"vpn-service/internal/service/vpn"
)

func RegisterRoutes(s *http.Server, service vpn.ServiceInterface) {
	registerApi(s.GetApiGroup(), service)
}

func registerApi(r *gin.RouterGroup, service vpn.ServiceInterface) {
	vpnGroup := r.Group("/vpn")
	{
		vpnGroup.POST("/register", RegisterUser(service))
		vpnGroup.GET("/by-chat/:chatId", GetUserByChatId(service))
		vpnGroup.POST("/reset-traffic/:chatId", ResetTraffic(service))
		vpnGroup.GET("/exists/:chatId", UserExists(service))
		vpnGroup.GET("/by-email/:email", GetUserByEmail(service))
		vpnGroup.GET("/subscription-link/:chatId", GetSubcLinkByChatId(service))
		vpnGroup.PUT("/by-uuid/:uuid/update-client", UpdateClient(service))
	}
}
