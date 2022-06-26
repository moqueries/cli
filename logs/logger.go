package logs

import (
	"log"
	"os"
)

var (
	debug       = true
	debugLogger *log.Logger
	warnLogger  *log.Logger
	errLogger   *log.Logger
)

func init() {
	debugLogger = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)
	warnLogger = log.New(os.Stderr, "WARNING: ", log.LstdFlags)
	errLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)
	log.SetPrefix("INFO: ")
}

// Init initializes the loggers
func Init(dbg bool) {
	debug = dbg
}

// IsDebug returns true if debug logging is enabled
func IsDebug() bool {
	return debug
}

// Debugf will output a debugging log if debug logs are configured
func Debugf(format string, args ...interface{}) {
	if debug {
		debugLogger.Printf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Warn logs a warning message
func Warn(msg string) {
	warnLogger.Print(msg)
}

// Error logs a error message with an error
func Error(msg string, err error) {
	errLogger.Printf(msg+": %v", err)
}

// Errorf logs a formatted message
func Errorf(format string, args ...interface{}) {
	errLogger.Printf(format, args...)
}

// Panic logs a message with an error then panics
func Panic(msg string, err error) {
	log.Panicf(msg+": %v", err)
}

// Panicf logs a formatted message then panics
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}
