//go:build !linux

package systemd

func Notify(state string) (bool, error) {
	return false, ErrUnsupported
}
