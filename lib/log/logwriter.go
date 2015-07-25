// Copyright 2014 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package log

import (
	. "github.com/limetext/lime-backend/lib/util"
	"github.com/limetext/log4go"
	"sync"
)

type LogWriter interface {
	log4go.LogWriter
}

func NewConsoleLogWriter() LogWriter {
	return log4go.NewConsoleLogWriter()
}

func NewFileLogWriter(fname string, rotate bool) LogWriter {
	return log4go.NewFileLogWriter(fname, rotate)
}

// Implementation of a default LogWriter which takes a handler function

type logWriter struct {
	sync.Mutex
	log chan string
}

func NewLogWriter(h func(string)) LogWriter {
	l := &logWriter{
		log: make(chan string, 100),
	}
	go func() {
		for fl := range l.log {
			h(fl)
		}
	}()
	return l
}

func (l *logWriter) LogWrite(rec *log4go.LogRecord) {
	p := Prof.Enter("log")
	defer p.Exit()
	l.Lock()
	defer l.Unlock()
	fl := log4go.FormatLogRecord(log4go.FORMAT_DEFAULT, rec)
	l.log <- fl
}

func (l *logWriter) Close() {
	l.Lock()
	defer l.Unlock()
	close(l.log)
}
