// gosnippets (c) 2023-2025 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

//go:build !linux

package systemd

func Notify(state string) (bool, error) {
	return false, ErrUnsupported
}
