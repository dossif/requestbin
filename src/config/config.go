package config

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	Listen    string
	Log       *logrus.Entry
	WaitGroup *sync.WaitGroup
}

var Service Config
