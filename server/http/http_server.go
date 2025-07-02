package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	serviceName = "todo_project-service"
	version = "0.0.1"
)

type Server struct {
	httpServer *http.Server
}

func New(port string, engine *gin.Engine) *Server {
	return &Server{
		httpServer: &http.Server {
			Addr: ":" + port,
			Handler: engine,
		},
	}
}

func NewEngine() *gin.Engine {
	engine := gin.New()
	engine.GET("/service", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"version": version,
		})
	})
	return engine
}

func (s *Server) Run() error {
	go func() {
		logrus.Infof("Service %s listening on %s", serviceName, s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	logrus.Info("Server exiting")
	return nil
}