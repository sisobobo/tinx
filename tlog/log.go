package tlog

import "github.com/golang/glog"

func Info(args ...any) {
	glog.V(3).Info(args...)
}

func Infof(format string, args ...any) {
	glog.V(3).Infof(format, args...)
}

func Warn(args ...any) {
	glog.V(2).Info(args...)
}

func Warnf(format string, args ...any) {
	glog.V(2).Infof(format, args...)
}

func Error(args ...any) {
	glog.V(1).Info(args...)
}

func Errorf(format string, args ...any) {
	glog.V(1).Infof(format, args...)
}
