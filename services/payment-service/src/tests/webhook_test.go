package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"payment-service/internal/controller/yoomoney"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestCreateYoomoneyLog(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing label", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		form := url.Values{}
		form.Add("notification_type", "incoming-transfer")
		form.Add("operation_id", "op123456")
		form.Add("amount", "100.50")
		form.Add("withdraw_amount", "99.00")
		form.Add("datetime", time.Now().Format(time.RFC3339))
		form.Add("label", time.Now().Format(time.RFC3339))
		// No label

		c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		yoomoney.CreateYoomoneyLog(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{}`, w.Body.String())
	})

	t.Run("valid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		form := url.Values{}
		form.Add("notification_type", "incoming-transfer")
		form.Add("operation_id", "op123456")
		form.Add("amount", "100.50")
		form.Add("withdraw_amount", "99.00")
		form.Add("datetime", time.Now().Format(time.RFC3339))
		form.Add("label", "payment-uuid")

		c.Request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		yoomoney.CreateYoomoneyLog(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{}`, w.Body.String())
	})
}
