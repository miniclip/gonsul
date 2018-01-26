package errorutil

import (
	"fmt"
	"time"
)

const LogErr 		= "ERROR"
const LogInfo 		= "INFO"
const LogDebug 		= "DEBUG"

const LogLevelErr 	= 10
const LogLevelInfo 	= 20
const LogLevelDebug	= 30

// Levels are the log levels we respond to=o.
var ErrorLevels = map[string]int {
	LogErr: 	LogLevelErr,
	LogInfo:  	LogLevelInfo,
	LogDebug: 	LogLevelDebug,
}

type Logger struct {
	level int
}

func NewLogger(level int) *Logger {
	return &Logger{
		level: level,
	}
}

func (logger *Logger) PrintError(msg string) {
	t := time.Now()
	fmt.Println("[" + LogErr + "] [" + t.Format(time.StampMilli) + "] " + msg)
}

func (logger *Logger) PrintInfo(msg string) {
	t := time.Now()
	if logger.level >= ErrorLevels[LogInfo] {
		fmt.Println("[" + LogInfo + "] [" + t.Format(time.StampMilli) + "] " + msg)
	}
}

func (logger *Logger) PrintDebug(msg string) {
	t := time.Now()
	if logger.level >= ErrorLevels[LogDebug] {
		fmt.Println("[" + LogDebug + "] [" + t.Format(time.StampMilli) + "] " + msg)
	}
}