package app

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Logger is a logger that has *os.File.
type Logger struct {
	*log.Logger
	file *os.File
}

func NewLogger(logfile string) (*Logger, error) {
	var file *os.File
	var w io.Writer

	if logfile == "" {
		w = os.Stderr
	} else if logfile == "-" {
		w = os.Stdout
	} else {
		path, err := absPath(logfile)
		if err != nil {
			return nil, err
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		w = file
	}

	logger := &Logger{
		Logger: log.New(w, "", log.LstdFlags),
		file:   file,
	}
	return logger, nil
}

func NewNullLogger() (*Logger, error) {
	return &Logger{
		Logger: log.New(ioutil.Discard, "", 0),
	}, nil
}

// Close close the log file when it is not nil.
func (l *Logger) Close() error {
	if l.file != nil {
		err := l.file.Close()
		if err != nil {
			return err
		}
		l.file = nil
	}
	return nil
}
