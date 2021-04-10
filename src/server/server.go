package server

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"requestbin/src/config"
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
	log.Infof("start http")
	mw := middleware{Service: service}
	h := handler{Service: service}
	router := mux.NewRouter()
	router.Use(mw.commonMiddleware)
	router.HandleFunc("/favicon.ico", h.handlerFavicon)
	router.HandleFunc("/", h.handlerRequestStatus)
	router.HandleFunc("/{status}", h.handlerRequestStatus)
	http.Handle("/", router)
	server := &http.Server{
		Addr:         service.Listen,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 5,
		Handler:      router,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start http server: %v", err)
		}
	}()
	<-osSignalChan
	ctx, cancel := context.WithTimeout(context.Background(), exitTimeout)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Errorf("failed to shutdown http server")
	}
	defer log.Infof("stop http")
}

// For common requests
func (s *middleware) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		log := s.Service.Log
		logFields := logrus.Fields{"remote": r.RemoteAddr, "method": r.Method, "uri": r.RequestURI, "host": r.Host}
		log.WithFields(logFields).Info("http-request")
		next.ServeHTTP(w, r)
	})
}
