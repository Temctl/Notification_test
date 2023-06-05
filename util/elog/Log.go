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
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	workerName := os.Getenv("workername")
	if workerName == "" {
		workerName = "default"
	}
	logPath := os.Getenv("logpath")
	if logPath == "" {
		logPath = "./logs/" + workerName + ".log"
	} else {
		logPath = logPath + workerName + ".log"
	}
	// ------------------------------------------------------
	// IMPORT ENV -------------------------------------------
	// ------------------------------------------------------
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
	InfoLogger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    elog_maxsize, // Max size in megabytes before log rotation
		MaxBackups: elog_backups, // Max number of old log files to keep
		MaxAge:     elog_maxage,  // Max number of days to keep log files
		Compress:   true,
	}, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
	WarnLogger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    elog_maxsize,
		MaxBackups: elog_backups,
		MaxAge:     elog_maxage,
		Compress:   true,
	}, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
	ErrorLogger = log.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    elog_maxsize,
		MaxBackups: elog_backups,
		MaxAge:     elog_maxage,
		Compress:   true,
	}, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
	// -------------------------------------------------------
	// SETCONFIG ---------------------------------------------
	// -------------------------------------------------------
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// -------------------------------------------------------
// PUBLIC FUNC -------------------------------------------
// -------------------------------------------------------
func Info() *log.Logger {
	return InfoLogger
}
func Warning() *log.Logger {
	return WarnLogger
}
func Error() *log.Logger {
	return ErrorLogger
}
