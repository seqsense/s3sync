package s3sync

import "log"

// LoggerIF is the logger interface which this library requires.
type LoggerIF interface {
	// Log inserts a log entry. Arguments are handled in the manner
	// of fmt.Print.
	Log(v ...interface{})
	// Log inserts a log entry. Arguments are handled in the manner
	// of fmt.Printf.
	Logf(format string, v ...interface{})
}

// Logger is the logger instance.
var Logger LoggerIF

// SetLogger Sets the logger.
func SetLogger(logger LoggerIF) {
	Logger = logger
}

func println(v ...interface{}) {
	if Logger == nil {
		log.Println(v...)
		return
	}
	Logger.Log(v...)
}
