// gosnippets (c) 2023-2025 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package systemd

import "errors"

const (
	Ready     = "READY=1"
	Stopping  = "STOPPING=1"
	Reloading = "RELOADING=1"
	Watchdog  = "WATCHDOG=1"
)

var (
	ErrUnsupported = errors.New("systemd is not supported on current platform")
)
