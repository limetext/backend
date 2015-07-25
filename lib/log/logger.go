// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	"github.com/limetext/log4go"
)

type Logger interface {
	AddFilter(name string, level Level, writer LogWriter)
	Finest(arg0 interface{}, args ...interface{})
	Fine(arg0 interface{}, args ...interface{})
	Debug(arg0 interface{}, args ...interface{})
	Trace(arg0 interface{}, args ...interface{})
	Info(arg0 interface{}, args ...interface{})
	Warn(arg0 interface{}, args ...interface{}) error
	Error(arg0 interface{}, args ...interface{}) error
	Errorf(format string, args ...interface{})
	Critical(arg0 interface{}, args ...interface{}) error
	Logf(level Level, format string, args ...interface{})
	Close()
}

type logger struct {
	log4go.Logger
}

func NewLogger() Logger {
	return &logger{make(log4go.Logger)}
}

func (l *logger) AddFilter(name string, level Level, writer LogWriter) {
	lvl := log4go.INFO
	switch level {
	case FINEST:
		lvl = log4go.FINEST
	case FINE:
		lvl = log4go.FINE
	case DEBUG:
		lvl = log4go.DEBUG
	case TRACE:
		lvl = log4go.TRACE
	case INFO:
		lvl = log4go.INFO
	case WARNING:
		lvl = log4go.WARNING
	case ERROR:
		lvl = log4go.ERROR
	case CRITICAL:
		lvl = log4go.CRITICAL
	default:
	}
	l.Logger.AddFilter(name, lvl, writer)
}

func (l *logger) Logf(level Level, format string, args ...interface{}) {
	lvl := log4go.INFO
	switch level {
	case FINEST:
		lvl = log4go.FINEST
	case FINE:
		lvl = log4go.FINE
	case DEBUG:
		lvl = log4go.DEBUG
	case TRACE:
		lvl = log4go.TRACE
	case INFO:
		lvl = log4go.INFO
	case WARNING:
		lvl = log4go.WARNING
	case ERROR:
		lvl = log4go.ERROR
	case CRITICAL:
		lvl = log4go.CRITICAL
	default:
	}
	l.Logger.Logf(lvl, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logf(ERROR, format, args...)
}
