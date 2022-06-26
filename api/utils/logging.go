package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func LogInfo(msg string, sprintfArgs ...any) {
	conf := GetConf()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f, err := os.OpenFile(conf.ApiLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		l := log.New(os.Stderr, "", 0)
		l.Printf("Could not log to %s, error: %s\n", conf.ApiLogFile, err.Error())
	}
	defer f.Close()

	log_file_logger := log.New(f, timestamp+" ", 0)
	log_file_logger.Printf("INFO: %s\n", fmt.Sprintf(msg, sprintfArgs...))
	std_out_logger := log.New(os.Stdout, timestamp+" ", 0)
	std_out_logger.Printf("%s\n", msg)
}

func LogDebug(msg string, sprintfArgs ...any) {
	conf := GetConf()
	if !conf.DebugMode {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f, err := os.OpenFile(conf.ApiLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		l := log.New(os.Stderr, "", 0)
		l.Printf("Could not log to %s, error: %s\n", conf.ApiLogFile, err.Error())
	}
	defer f.Close()

	log_file_logger := log.New(f, timestamp+" ", 0)
	log_file_logger.Printf("WARNING: %s\n", fmt.Sprintf(msg, sprintfArgs...))
	std_out_logger := log.New(os.Stdout, timestamp+" ", 0)
	std_out_logger.Printf("%s\n", msg)
}

func LogError(msg string, sprintfArgs ...any) {
	conf := GetConf()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f, err := os.OpenFile(conf.ApiLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		l := log.New(os.Stderr, "", 0)
		l.Printf("Could not log to %s, error: %s\n", conf.ApiLogFile, err.Error())
	}
	defer f.Close()

	log_file_logger := log.New(f, timestamp+" ", 0)
	log_file_logger.Printf("ERROR: %s\n", fmt.Sprintf(msg, sprintfArgs...))
	std_err_logger := log.New(os.Stderr, timestamp+" ", 0)
	std_err_logger.Printf("%s\n", msg)
}
