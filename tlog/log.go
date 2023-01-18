package tlog

import "log"

func Infof(format string, msg ...any) {
	log.Printf(format, msg...)
}

func Warnf(format string, msg ...any) {
	log.Printf(format, msg...)
}

func Errorf(format string, msg ...any) {
	log.Printf(format, msg...)
}
