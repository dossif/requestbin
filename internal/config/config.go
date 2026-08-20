package config

import (
	"log/slog"
	"sync"
)

type Config struct {
	Listen    string
	Log       *slog.Logger
	WaitGroup *sync.WaitGroup
}

var Service Config
