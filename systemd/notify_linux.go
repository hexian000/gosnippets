// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

//go:build linux

package systemd

import (
	"net"
	"os"
)

// Notify sends a notification message to the systemd init system.
// It returns true if the notification was sent.
func Notify(state string) (bool, error) {
	addr := os.Getenv("NOTIFY_SOCKET")
	if addr == "" {
		return false, nil
	}

	conn, err := net.Dial("unixgram", addr)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(state)); err != nil {
		return false, err
	}
	return true, nil
}
