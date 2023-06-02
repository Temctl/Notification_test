package elog

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/natefinch/lumberjack.v2"
)

// -------------------------------------------------------
// PUBLIC VAR --------------------------------------------
// -------------------------------------------------------
var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	workerName := os.Getenv("workername")
	if workerName == "" {
		workerName = "default"
	}
	// ------------------------------------------------------
	// IMPORT ENV -------------------------------------------
	// ------------------------------------------------------
	logPath := os.Getenv("ELOG_PATH")
	if logPath == "" {
		logPath = "./logs/" + workerName + ".log"
	}
	elog_maxsize, err := strconv.Atoi(os.Getenv("ELOG_MAXSIZE"))
	if err != nil {
		elog_maxsize = 100
	}
	elog_backups, err := strconv.Atoi(os.Getenv("ELOG_BACKUPS"))
	if err != nil {
		elog_backups = 10
	}
	elog_maxage, err := strconv.Atoi(os.Getenv("ELOG_MAXAGE"))
	if err != nil {
		elog_maxage = 30
	}
	// -------------------------------------------------------
	// CREATE LOGGER -----------------------------------------
	// -------------------------------------------------------
	logFile := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    elog_maxsize, // Max size in megabytes before log rotation (e.g., 10MB)
		MaxBackups: elog_backups, // Max number of old log files to keep
		MaxAge:     elog_maxage,  // Max number of days to keep old log files
		Compress:   true,         // Whether to compress old log files
	}
	infoLogger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    elog_maxsize, // Max size in megabytes before log rotation
		MaxBackups: elog_backups, // Max number of old log files to keep
		MaxAge:     elog_maxage,  // Max number of days to keep log files
		Compress:   true,
	}, "INFO: ", log.Ldate|log.Ltime)
	warnLogger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    elog_maxsize,
		MaxBackups: elog_backups,
		MaxAge:     elog_maxage,
		Compress:   true,
	}, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    elog_maxsize,
		MaxBackups: elog_backups,
		MaxAge:     elog_maxage,
		Compress:   true,
	}, "ERROR: ", log.Ldate|log.Ltime)
	// -------------------------------------------------------
	// SETCONFIG ---------------------------------------------
	// -------------------------------------------------------
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// -------------------------------------------------------
// PUBLIC FUNC -------------------------------------------
// -------------------------------------------------------
func Info(message string) {
	infoLogger.Println(message)
}
func Warning(message error) {
	warnLogger.Println(message)
}
func Error(message string, err error) {
	errorLogger.Printf(" %s : %s", message, err)
}
