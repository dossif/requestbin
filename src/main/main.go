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
func initLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "timestamp",
		logrus.FieldKeyLevel: "level",
		logrus.FieldKeyMsg:   "message"}})
	logger.SetReportCaller(false)
	switch logLevel {
	case "DEBUG":
		logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		logger.SetLevel(logrus.InfoLevel)
	case "ERROR":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		panic(fmt.Sprintf("log level %v is not (DEBUG|INFO|ERROR)", logLevel))
	}
	return logger
}

// Main
func main() {
	fmt.Println(appLogo)
	listen, ok := os.LookupEnv("RB_LISTEN")
	if ok != true {
		listen = "0.0.0.0:8080"
	}
	logLevel, ok := os.LookupEnv("RB_LOGLEVEL")
	if ok != true {
		logLevel = "ERROR"
	}
	log := initLogger(logLevel).WithField("app", strings.ToLower(appLogo))
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
	fmt.Println(fmt.Sprintf("loglevel: %v", logLevel))
	go server.Server(osSignalChan)
	waitGroup.Add(1)
	waitGroup.Wait()
	log.Info("stop app")
	os.Exit(0)
}
