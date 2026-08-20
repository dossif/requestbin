package main

import (
	"os"
	"strings"
	"sync"

	"github.com/dossif/requestbin/internal/config"
	"github.com/dossif/requestbin/internal/server"
	"github.com/dossif/requestbin/pkg/logger"
	"github.com/dossif/requestbin/pkg/shutdown"
)

const appName = "requestbin"

var AppVersion = "0.0.0"

// Main
func main() {
	listen, ok := os.LookupEnv("RB_LISTEN")
	if ok != true {
		listen = "0.0.0.0:8080"
	}
	log := logger.New().With("app", strings.ToLower(appName))
	osSignalChan := make(chan bool)
	shutdown.WaitForSignal(log, osSignalChan)
	var waitGroup sync.WaitGroup
	service := config.Config{
		Listen:    listen,
		WaitGroup: &waitGroup,
		Log:       log,
	}
	config.Service = service
	log.Info("start", "version", AppVersion, "listen", listen)
	go server.Server(osSignalChan)
	waitGroup.Add(1)
	waitGroup.Wait()
	log.Info("stop")
	os.Exit(0)
}
