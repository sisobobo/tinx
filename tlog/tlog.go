package tlog

import "log"

func INFO(format string, params ...interface{}) {
	log.Printf(format, params...)
}

func WARN(format string, params ...interface{}) {
	log.Printf(format, params...)
}

func Error(format string, params ...interface{}) {
	log.Printf(format, params...)
}
