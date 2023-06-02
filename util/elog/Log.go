package elog

import (
	"fmt"
	"log"
	"os"
	"reflect"
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
	// ------------------------------------------------------
	// IMPORT ENV -------------------------------------------
	// ------------------------------------------------------
	ELOG_PATH := os.Getenv("ELOG_PATH")
	if ELOG_PATH == "" {
		ELOG_PATH = "./tmp"
	}
	ELOG_MAXSIZE := os.Getenv("ELOG_MAXSIZE")
	if ELOG_MAXSIZE == "" {
		ELOG_MAXSIZE = "100"
	}
	ELOG_BACKUPS := os.Getenv("ELOG_BACKUPS")
	if ELOG_BACKUPS == "" {
		ELOG_BACKUPS = "10"
	}
	ELOG_MAXAGE := os.Getenv("ELOG_MAXAGE")
	if ELOG_MAXAGE == "" {
		ELOG_MAXAGE = "30"
	}
	// ------------------------------------------------------
	// CONVERT TYPE -----------------------------------------
	// ------------------------------------------------------
	elog_maxsize, err := strconv.Atoi(ELOG_MAXSIZE)
	fmt.Println(elog_maxsize, err, reflect.TypeOf(elog_maxsize))
	elog_backups, err := strconv.Atoi(ELOG_BACKUPS)
	fmt.Println(elog_backups, err, reflect.TypeOf(elog_backups))
	elog_maxage, err := strconv.Atoi(ELOG_MAXAGE)
	fmt.Println(elog_maxage, err, reflect.TypeOf(elog_maxage))
	// -------------------------------------------------------
	// CREATE LOGGER -----------------------------------------
	// -------------------------------------------------------
	logFile := &lumberjack.Logger{
		Filename:   ELOG_PATH + "/inputWorker/inputWorker.log",
		MaxSize:    elog_maxsize, // Max size in megabytes before log rotation (e.g., 10MB)
		MaxBackups: elog_backups, // Max number of old log files to keep
		MaxAge:     elog_maxage,  // Max number of days to keep old log files
		Compress:   true,         // Whether to compress old log files
	}
	infoLogger = log.New(&lumberjack.Logger{
		Filename:   ELOG_PATH + "/info/info.log",
		MaxSize:    elog_maxsize, // Max size in megabytes before log rotation
		MaxBackups: elog_backups, // Max number of old log files to keep
		MaxAge:     elog_maxage,  // Max number of days to keep log files
		Compress:   true,
	}, "INFO: ", log.Ldate|log.Ltime)
	warnLogger = log.New(&lumberjack.Logger{
		Filename:   ELOG_PATH + "/warning/warning.log",
		MaxSize:    elog_maxsize,
		MaxBackups: elog_backups,
		MaxAge:     elog_maxage,
		Compress:   true,
	}, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(&lumberjack.Logger{
		Filename:   ELOG_PATH + "/error/error.log",
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
