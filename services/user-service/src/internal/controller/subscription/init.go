package subscription

import (
	"github.com/gin-gonic/gin"
	"pkg/infrastructure/http"
	"user-service/internal/repository/user/sqlite"
)

func RegisterRoutes(s *http.Server, userRepo *sqlite.UserRepository) {
	registerApi(s.GetApiGroup(), userRepo)
}

func registerApi(r *gin.RouterGroup, userRepo *sqlite.UserRepository) {

}
