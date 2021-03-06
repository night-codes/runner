package main

import (
	"strings"
	"time"
)

const (
	typeInfo    = 0
	typeTitle   = 1
	typeError   = 2
	typeWarning = 3
)

type (
	logger struct {
		Type    int
		Service *serviceStruct
		buffer  string
	}
)

func (l *logger) Write(b []byte) (n int, err error) {
	return l.WriteString(string(b))
}

func (l *logger) WriteString(s string) (n int, err error) {
	l.buffer = l.buffer + s
	split := strings.Split(l.buffer, "\n")

	if len(split) > 1 {
		l.buffer = split[len(split)-1]
		for _, v := range split[:len(split)-1] {
			if l.Service.Status == statusWaiting && (l.Service.StartMessage == "" || strings.Contains(v, l.Service.StartMessage)) {
				l.Service.changeStatus(statusRunned)
			}
			if m := strings.TrimSpace(v); v != "" {
				msg := logMessage{Time: time.Now(), Message: m, Type: l.Type}
				mainWS.Send("log", obj{"end": len(l.Service.Logs), "message": msg, "service": l.Service.ID})
				l.Service.Logs = append(l.Service.Logs, msg)
			}
		}
	}

	return len(s), nil
}
