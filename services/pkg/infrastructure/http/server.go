package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/novikoff-vvs/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	r        *gin.Engine
	log      logger.Interface
	apiGroup *gin.RouterGroup
	webGroup *gin.RouterGroup
}

func (s *Server) GetApiGroup() *gin.RouterGroup {
	return s.apiGroup
}

func (s *Server) GetWebGroup() *gin.RouterGroup {
	return s.webGroup
}

func (s *Server) Run(port string) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.r,
	}
	var errCh = make(chan error, 1)
	go func(errCh chan error) {
		log.Printf("Server is starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}(errCh)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errCh:
		{
			return err
		}
	case <-quit:
		{
			<-quit
			log.Println("Shutting down server...")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server exiting")
	return nil
}
func (s *Server) RegisterStatic() {
	s.r.LoadHTMLGlob("templates/*")
}

func NewServer(log logger.Interface) *Server {
	r := gin.Default()

	apiGroup := r.Group("/api")
	webGroup := r.Group("/web")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://t.me/nvs_vpn_3x_bot")
	})
	return &Server{
		r:        r,
		log:      log,
		apiGroup: apiGroup,
		webGroup: webGroup,
	}
}
