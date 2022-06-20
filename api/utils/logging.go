package utils

import (
	"log"
	"os"
)

func LogInfo(msg string) {
	conf := GetConf()
	f, err := os.OpenFile(conf.ApiLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		l := log.New(os.Stderr, "", 0)
		l.Printf("Could not log to %s, error: %s\n", conf.ApiLogFile, err.Error())
	}
	defer f.Close()

	log_file_logger := log.New(f, "", 0)
	log_file_logger.Printf("INFO: %s\n", msg)
	std_out_logger := log.New(os.Stdout, "", 0)
	std_out_logger.Printf("%s\n", msg)
}

func LogError(msg string) {
	conf := GetConf()
	f, err := os.OpenFile(conf.ApiLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		l := log.New(os.Stderr, "", 0)
		l.Printf("Could not log to %s, error: %s\n", conf.ApiLogFile, err.Error())
	}
	defer f.Close()

	log_file_logger := log.New(f, "", 0)
	log_file_logger.Printf("ERROR: %s\n", msg)
	std_err_logger := log.New(os.Stderr, "", 0)
	std_err_logger.Printf("%s\n", msg)
}
