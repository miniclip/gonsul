package errorutil

import (
	"fmt"
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
	fmt.Println("[" + LogErr + "] " + msg)
}

func (logger *Logger) PrintInfo(msg string) {
	if logger.level >= ErrorLevels[LogInfo] {
		fmt.Println("[" + LogInfo + "]  " + msg)
	}
}

func (logger *Logger) PrintDebug(msg string) {
	if logger.level >= ErrorLevels[LogDebug] {
		fmt.Println("[" + LogDebug + "] " + msg)
	}
}