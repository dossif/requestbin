package server

import (
	"context"
	"github.com/dossif/requestbin/internal/config"
	"net/http"
	"time"
)

const (
	exitTimeout = 5 * time.Second
)

type middleware struct {
	Service config.Config
}

type handler struct {
	Service config.Config
}

// Http Server
func Server(osSignalChan chan bool) {
	service := config.Service
	defer service.WaitGroup.Done()
	log := service.Log
	log.Info("start http server")
	mw := middleware{Service: service}
	h := handler{Service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", h.handlerFavicon)
	mux.HandleFunc("/status/{status}", h.handlerRequestStatus)
	mux.HandleFunc("/", h.handlerRequestStatus)
	router := mw.commonMiddleware(mux)
	server := &http.Server{
		Addr:         service.Listen,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 5,
		Handler:      router,
	}
	serveErrChan := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErrChan <- err
		}
	}()
	select {
	case <-osSignalChan:
	case err := <-serveErrChan:
		log.Error("http server failed to start", "error", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), exitTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown http server", "error", err)
	}
	log.Info("stop http server")
}

// For common requests
func (s *middleware) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		log := s.Service.Log
		log.Info("http-request", "remote", r.RemoteAddr, "method", r.Method, "uri", r.RequestURI, "host", r.Host)
		next.ServeHTTP(w, r)
	})
}
