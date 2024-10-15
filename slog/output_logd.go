//go:build android
// +build android

package slog

import (
	"encoding/binary"
	"net"
	"os"
	"strconv"
)

type logdWriter struct {
	tag []byte
	out net.Conn
}

func init() {
	newSyslogWriter = func(tag string) output {
		conn, err := net.Dial("unixgram", "/dev/socket/logdw")
		if err != nil {
			panic(err)
		}
		return &logdWriter{
			tag: []byte(tag),
			out: conn,
		}
	}
}

var levelMap = [...]byte{
	LevelSilence:     8, /* ANDROID_LOG_SILENT */
	LevelFatal:       7, /* ANDROID_LOG_FATAL */
	LevelError:       6, /* ANDROID_LOG_ERROR */
	LevelWarning:     5, /* ANDROID_LOG_WARN */
	LevelNotice:      4, /* ANDROID_LOG_INFO */
	LevelInfo:        4, /* ANDROID_LOG_INFO */
	LevelDebug:       3, /* ANDROID_LOG_DEBUG */
	LevelVerbose:     2, /* ANDROID_LOG_VERBOSE */
	LevelVeryVerbose: 2, /* ANDROID_LOG_VERBOSE */
}

func (l *logdWriter) Write(m *message) error {
	buf := make([]byte, 11, bufSize)
	buf[0] = 0 // LOG_ID_MAIN
	le := binary.LittleEndian
	le.PutUint16(buf[1:3], uint16(os.Getpid()))
	now := m.timestamp.UnixNano()
	le.PutUint32(buf[3:7], uint32(now/1000000000))
	le.PutUint32(buf[7:11], uint32(now%1000000000))

	buf = append(buf, levelMap[m.level])
	buf = append(buf, l.tag...)
	buf = append(buf, 0)

	buf = append(buf, levelChar[m.level], ' ')
	buf = append(buf, m.file...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(m.line), 10)
	buf = append(buf, ' ')
	buf = m.appendMessage(buf)
	buf = append(buf, 0)
	_, err := l.out.Write(buf)
	return err
}
