package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger = &logrus.Logger{
	Out:   os.Stdout,
	Hooks: nil,
	Formatter: &logrus.TextFormatter{
		ForceColors:               true,
		ForceQuote:                true,
		EnvironmentOverrideColors: true,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           time.RFC1123Z,
		DisableSorting:            true,
		PadLevelText:              true,
		QuoteEmptyFields:          true,
	},
	ReportCaller: false,
	Level:        logrus.InfoLevel,
}
