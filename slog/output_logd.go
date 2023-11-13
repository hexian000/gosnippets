//go:build android
// +build android

package slog

import (
	"encoding/binary"
	"net"
	"os"
	"strconv"
	"sync"
)

type logdOutput struct {
	mu  sync.Mutex
	tag []byte
	buf []byte
	out net.Conn
}

func init() {
	builtinOutput["syslog"] = func(tag string) (logOutput, error) {
		conn, err := net.Dial("unixgram", "/dev/socket/logdw")
		if err != nil {
			return nil, err
		}
		return &logdOutput{
			tag: []byte(tag),
			buf: make([]byte, 11), // android_log_header_t
			out: conn,
		}, nil
	}
}

var levelMap = [...]byte{
	LevelSilence: 8, /* ANDROID_LOG_SILENT */
	LevelFatal:   7, /* ANDROID_LOG_FATAL */
	LevelError:   6, /* ANDROID_LOG_ERROR */
	LevelWarning: 5, /* ANDROID_LOG_WARN */
	LevelNotice:  4, /* ANDROID_LOG_INFO */
	LevelInfo:    4, /* ANDROID_LOG_INFO */
	LevelDebug:   3, /* ANDROID_LOG_DEBUG */
	LevelVerbose: 2, /* ANDROID_LOG_VERBOSE */
}

func (l *logdOutput) Write(m logMessage) {
	l.mu.Lock()
	defer l.mu.Unlock()
	buf := l.buf[:11]
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
	buf = append(buf, m.msg...)
	buf = append(buf, 0)
	l.buf = buf
	_, _ = l.out.Write(buf)
}
