package utils

import (
	"log"
	"os"
	"time"

	"github.com/jimlawless/whereami"
)

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
)

var logLevelMap map[string]int = map[string]int{
	"DEBUG": LOG_LEVEL_DEBUG,
	"INFO":  LOG_LEVEL_INFO,
	"WARN":  LOG_LEVEL_WARN,
	"ERROR": LOG_LEVEL_ERROR,
}

func DebugLog(args ...interface{}) {
	printLog("DEBUG", args...)
}

func InfoLog(args ...interface{}) {
	printLog("INFO", args...)
}

func WarnLog(args ...interface{}) {
	printLog("WARN", args...)
}

func ErrorLog(args ...interface{}) {
	printLog("ERROR", args...)
}

func printLog(logLevel string, args ...interface{}) {
	appLogLevel, ok := os.LookupEnv("APP_LOG_LEVEL")
	if !ok || IsEmptyString(appLogLevel) {
		appLogLevel = "WARN"
	}
	if appLvl, ok := logLevelMap[appLogLevel]; ok && appLvl > logLevelMap[logLevel] {
		return
	}
	currentTime := time.Now()

	os.Mkdir("storage/logs", 0755)
	path := "storage/logs/dropbox-"
	fileName := path + currentTime.Format("2006-01-02") + ".log"

	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, logLevel+": ", log.Ldate|log.Ltime|log.Lshortfile)
	args = append(args, whereami.WhereAmI(3))
	logger.Println(args...)
}
