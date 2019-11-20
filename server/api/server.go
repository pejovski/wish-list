package api

import (
	"context"
	"fmt"
	"github.com/pejovski/wish-list/controller"
	srv "github.com/pejovski/wish-list/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const (
	ReadTimeout  = time.Second * 3
	WriteTimeout = time.Second * 3
)

type server struct {
	router Router
}

func NewServer(c controller.Controller) srv.Server {
	return server{router: newRouter(c)}
}

func (s server) Run(ctx context.Context) {
	server := &http.Server{
		Handler:      s.router,
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	doneCh := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			logrus.Info("API server is shutting down")
			shutdownCtx, cancel := context.WithTimeout(
				context.Background(),
				time.Second*5,
			)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				logrus.Errorf("API Server error: %s", err)
			}
		case <-doneCh:
		}
	}()

	logrus.Infof("API Server started at port: %s", os.Getenv("APP_PORT"))
	if err := server.ListenAndServe(); err != nil {
		logrus.Errorf("API Server error: %s", err)
	}

	close(doneCh)
}
