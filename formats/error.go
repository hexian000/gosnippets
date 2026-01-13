// gosnippets (c) 2023-2026 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package formats

import (
	"errors"
	"fmt"
	"reflect"
)

var errFormat = errors.New("format error")

// Error returns a detailed string representation of the given error.
func Error(err error) string {
	if err == nil {
		return "nil"
	}
	// errors.errorString
	if reflect.TypeOf(err) == reflect.TypeOf(errFormat) {
		return err.Error()
	}
	return fmt.Sprintf("(%T) %v", err, err)
}
