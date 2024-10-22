// gosnippets (c) 2023-2024 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package formats

import (
	"errors"
	"fmt"
	"reflect"
)

var errFormat = errors.New("format error")

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
