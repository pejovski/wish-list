package factory

import (
	"net/http"
	"time"
)

const (
	serverReadTimeout  = 3 * time.Second
	serverWriteTimeout = 3 * time.Second
)

func CreateHttpServer(h http.Handler, addr string) *http.Server {
	return &http.Server{
		Handler:      h,
		Addr:         addr,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
	}
}
