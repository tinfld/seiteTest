package main

import (
	"log"
	"os"
)

type LogFile struct {
	LogUsers  *log.Logger
	LogMovie  *log.Logger
	LogEnv    *log.Logger
	LogMiddle *log.Logger
	LogMain   *log.Logger

	file []*os.File
}

func loadLog() (LogFile, error) {
	var logFile LogFile
	var fileName = [...]string{
		"log/LogUsers",
		"log/LogMovie",
		"log/LogEnv",
		"log/LogMiddle",
		"log/LogMain",
	}

	for _, name := range fileName {
		file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return LogFile{}, err
		}

		logger := log.New(file, "", log.Ldate|log.Ltime)

		if name == "log/LogUsers" {
			logFile.LogUsers = logger
		} else if name == "log/LogMovie" {
			logFile.LogMovie = logger
		} else if name == "log/LogEnv" {
			logFile.LogEnv = logger
		} else if name == "log/LogMiddle" {
			logFile.LogMiddle = logger
		} else if name == "log/LogMain" {
			logFile.LogMain = logger
		}
		logFile.file = append(logFile.file, file)
	}
	return logFile, nil
}

func (l LogFile) close() {
	for _, file := range l.file {
		file.Close()
	}
}
