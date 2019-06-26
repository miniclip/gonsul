package util

import (
	"fmt"
	"time"
)

const LogErr = "ERROR"
const LogInfo = "INFO"
const LogDebug = "DEBUG"

const LogLevelErr = 10
const LogLevelInfo = 20
const LogLevelDebug = 30

// Levels are the log levels we respond to=o.
var ErrorLevels = map[string]int{
	LogErr:   LogLevelErr,
	LogInfo:  LogLevelInfo,
	LogDebug: LogLevelDebug,
}

// ILogger is our logger interface
type ILogger interface {
	PrintError(msg string)
	PrintInfo(msg string)
	PrintDebug(msg string)
	AddMessage(msg string)
	GetMessages() []string
}

// logger is our ILogger interface concrete implementation. It's used throughout the
// application to print/output messages to stdout
type logger struct {
	level 		int
	messages 	[]string
}

// NewLogger is our logger constructor
func NewLogger(level int) ILogger {
	return &logger{
		level: level,
		messages: []string{},
	}
}

// PrintError prints a message of ERROR level
func (logger *logger) PrintError(msg string) {
	t := time.Now()
	fmt.Println("[" + LogErr + "] [" + t.Format(time.StampMilli) + "] " + msg)
}

// PrintInfo prints a message of INFO level
func (logger *logger) PrintInfo(msg string) {
	t := time.Now()
	if logger.level >= ErrorLevels[LogInfo] {
		fmt.Println("[" + LogInfo + "]  [" + t.Format(time.StampMilli) + "] " + msg)
	}
}

// PrintDebug prints a message of DEBUG level
func (logger *logger) PrintDebug(msg string) {
	t := time.Now()
	if logger.level >= ErrorLevels[LogDebug] {
		fmt.Println("[" + LogDebug + "] [" + t.Format(time.StampMilli) + "] " + msg)
	}
}

// AddMessage add a message to internal logger message slice
func (logger *logger) AddMessage(msg string) {
	logger.messages = append(logger.messages, msg)
}

// GetMessages returns our internal message slice
func (logger *logger) GetMessages() []string {
	return logger.messages
}
