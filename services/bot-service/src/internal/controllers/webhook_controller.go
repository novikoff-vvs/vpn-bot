package controllers

import (
	"bot-service/internal/vpn"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebhookController struct {
	vpnClient vpn.ServiceInterface
}

func NewWebhookController(vpnClient vpn.ServiceInterface) *WebhookController {
	return &WebhookController{vpnClient: vpnClient}
}

func (wc WebhookController) SetupRoutes(r gin.IRouter) {
	g := r.Group("webhook")
	{
		g.POST("handle", wc.handle)
	}
}

func (wc WebhookController) handle(c *gin.Context) {
	var req struct {
		ChatId int64 `json:"chat_id"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = wc.vpnClient.ResetClientTraffic(req.ChatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
