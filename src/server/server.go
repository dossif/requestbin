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

func Server(osSignalChan chan bool) {
	service := config.Service
	defer service.WaitGroup.Done()
	log := service.Log.WithField("service", "rest")
	log.Infof("start service")
	mw := middleware{Service: service}
	h := handler{Service: service}
	// common router
	router := mux.NewRouter()
	router.Use(mw.commonMiddleware)
	router.HandleFunc("/", h.handlerRequestStatus)
	router.HandleFunc("/status/{status}", h.handlerRequestStatus)
	http.Handle("/", router)
	// create http server
	server := &http.Server{
		Addr:         service.Listen,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 5,
		Handler:      router,
	}
	go func() {
		_ = server.ListenAndServe()
	}()
	<-osSignalChan
	log.Infof("stop service")
	ctx, cancel := context.WithTimeout(context.Background(), exitTimeout)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Errorf("failed to shutdown server")
	}
	defer log.Infof("stop service")
}

// For common requests
func (s *middleware) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		log := s.Service.Log
		logFields := logrus.Fields{"remote": r.RemoteAddr, "method": r.Method, "uri": r.RequestURI, "service": "rest"}
		log.WithFields(logFields).Info("http-request")
		next.ServeHTTP(w, r)
	})
}
