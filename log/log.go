package log

import (
	"fmt"
	"os"

	"github.com/raunaqrox/assistme/config"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	// create log folder if not exists
	if config.ENV == "PROD" {
		if _, err := os.Stat(config.LogFolder); os.IsNotExist(err) {
			os.Mkdir(config.LogFolder, os.ModeDir)
		}
	}

	f, err := os.OpenFile(config.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("Could not setup Log, Error opening file: %v", err))
	}

	Log = &logrus.Logger{
		Out:       f,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

}
