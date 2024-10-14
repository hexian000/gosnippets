//go:build !windows && !android && !plan9
// +build !windows,!android,!plan9

package slog

import (
	"log/syslog"
	"strconv"
)

type syslogWriter struct {
	out *syslog.Writer
}

func init() {
	newSyslogWriter = func(tag string) writer {
		w, err := syslog.New(syslog.LOG_USER|syslog.LOG_NOTICE, tag)
		if err != nil {
			panic(err)
		}
		return &syslogWriter{w}
	}
}

var priorityMap = [...]func(*syslog.Writer, string) error{
	(*syslog.Writer).Alert,
	(*syslog.Writer).Crit,
	(*syslog.Writer).Err,
	(*syslog.Writer).Warning,
	(*syslog.Writer).Notice,
	(*syslog.Writer).Info,
	(*syslog.Writer).Debug,
	(*syslog.Writer).Debug,
	(*syslog.Writer).Debug,
}

func (s *syslogWriter) Write(m *message) {
	buf := make([]byte, 0, bufSize)
	buf = append(buf, levelChar[m.level], ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	buf = append(buf, m.msg...)
	_ = priorityMap[m.level](s.out, string(buf))
}
