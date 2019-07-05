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
var logger LoggerIF

// SetLogger sets the logger.
func SetLogger(l LoggerIF) {
	logger = l
}

func println(v ...interface{}) {
	if logger == nil {
		log.Println(v...)
		return
	}
	logger.Log(v...)
}
