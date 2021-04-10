package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"requestbin/src/config"
	"requestbin/src/server"
	"strings"
	"sync"
	"syscall"
)

const appLogo = "RequestBin"

// Handle OS signals
func setOsSignalHandler(log *logrus.Entry, osSignalChan chan bool) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.WithField("signal", sig).Infof("os-signal")
		close(osSignalChan)
	}()
}

// Initiate logging system
func initLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "timestamp",
		logrus.FieldKeyLevel: "level",
		logrus.FieldKeyMsg:   "message"}})
	logger.SetReportCaller(false)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

// Main
func main() {
	fmt.Println(appLogo)
	listen, ok := os.LookupEnv("RB_LISTEN")
	if ok != true {
		listen = "0.0.0.0:8080"
	}
	log := initLogger().WithField("app", strings.ToLower(appLogo))
	osSignalChan := make(chan bool)
	setOsSignalHandler(log, osSignalChan)
	var waitGroup sync.WaitGroup
	service := config.Config{
		Listen:    listen,
		WaitGroup: &waitGroup,
		Log:       log,
	}
	config.Service = service
	fmt.Println(fmt.Sprintf("listen: %v", service.Listen))
	go server.Server(osSignalChan)
	waitGroup.Add(1)
	waitGroup.Wait()
	log.Info("stop app")
	os.Exit(0)
}
