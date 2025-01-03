// gosnippets (c) 2023-2025 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package formats

import (
	"fmt"
	"math"
	"time"
)

func isNormal(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0) && f != 0.0
}

func formatAbnormal(f float64) string {
	if math.IsNaN(f) {
		if math.Signbit(f) {
			return "-nan"
		}
		return "nan"
	}
	if math.IsInf(f, 0) {
		if math.Signbit(f) {
			return "-inf"
		}
		return "inf"
	}
	if f == 0.0 {
		if math.Signbit(f) {
			return "-0"
		}
		return "0"
	}
	return fmt.Sprintf("%e", f)
}

var siPrefixPos = [...]string{
	"k", "M", "G", "T", "P", "E", "Z", "Y", "R", "Q",
}

var siPrefixNeg = [...]string{
	"m", "μ", "n", "p", "f", "a", "z", "y", "r", "q",
}

func SIPrefix(value float64) string {
	if !isNormal(value) {
		return formatAbnormal(value)
	}
	absvalue := math.Abs(value)
	if !(1e-30 < absvalue && absvalue < 1e+31) {
		return fmt.Sprintf("%.2e", value)
	}
	e := int(math.Floor(math.Log10(absvalue) / 3.0))
	if e == 0 {
		return fmt.Sprintf("%.3g", value)
	}
	if e < 0 {
		i := -e
		if i > len(siPrefixNeg) {
			i = len(siPrefixNeg)
		}
		v := value / math.Pow(10.0, -3.0*float64(i))
		prefix := siPrefixNeg[i-1]
		return fmt.Sprintf("%.3g%s", v, prefix)
	}
	i := e
	if i > len(siPrefixPos) {
		i = len(siPrefixPos)
	}
	v := value / math.Pow(10.0, 3.0*float64(i))
	prefix := siPrefixPos[i-1]
	return fmt.Sprintf("%.3g%s", v, prefix)
}

var iecUnits = [...]string{
	"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB",
}

func IECBytes(value float64) string {
	if !isNormal(value) {
		return formatAbnormal(value)
	}
	e := (int(math.Log2(math.Abs(value))) - 1) / 10
	if e < 0 {
		e = 0
	} else if e >= len(iecUnits) {
		e = len(iecUnits) - 1
	}
	v := math.Ldexp(value, e*-10)
	if e > 0 {
		if -10.0 < v && v < 10.0 {
			return fmt.Sprintf("%.2f%s", v, iecUnits[e])
		}
		if -100.0 < v && v < 100.0 {
			return fmt.Sprintf("%.1f%s", v, iecUnits[e])
		}
	}
	return fmt.Sprintf("%.0f%s", v, iecUnits[e])
}

// DurationSeconds formats the truncated duration
func DurationSeconds(d time.Duration) string {
	sign := ""
	if d < 0 {
		sign = "-"
	}
	value := int64(d / time.Second)
	if value < 0 {
		value = -value
	}
	var second, minute, hour uint8
	second, value = uint8(value%60), value/60
	minute, value = uint8(value%60), value/60
	hour, value = uint8(value%24), value/24
	day := value
	if day != 0 {
		return fmt.Sprintf("%s%dd%02d:%02d:%02d", sign, day, hour, minute, second)
	}
	if hour != 0 {
		return fmt.Sprintf("%s%d:%02d:%02d", sign, hour, minute, second)
	}
	return fmt.Sprintf("%s%d:%02d", sign, minute, second)
}

// DurationMillis formats the truncated duration
func DurationMillis(d time.Duration) string {
	sign := ""
	if d < 0 {
		sign = "-"
	}
	value := int64(d / time.Millisecond)
	if value < 0 {
		value = -value
	}
	var milli uint16
	var second, minute, hour uint8
	milli, value = uint16(value%1000), value/1000
	second, value = uint8(value%60), value/60
	minute, value = uint8(value%60), value/60
	hour, value = uint8(value%24), value/24
	day := value
	if day != 0 {
		return fmt.Sprintf("%s%dd%02d:%02d:%02d.%03d", sign, day, hour, minute, second, milli)
	}
	if hour != 0 {
		return fmt.Sprintf("%s%d:%02d:%02d.%03d", sign, hour, minute, second, milli)
	}
	return fmt.Sprintf("%s%d:%02d.%03d", sign, minute, second, milli)
}

// DurationNanos formats the precise duration
func DurationNanos(d time.Duration) string {
	value := int64(d / time.Nanosecond)
	sign := ""
	s := int64(1)
	if value < 0 {
		sign = "-"
		s = int64(-1)
	}
	var nano uint32
	var second, minute, hour uint8
	nano, value = uint32(value%1000000000*s), value/1000000000*s
	second, value = uint8(value%60), value/60
	minute, value = uint8(value%60), value/60
	hour, value = uint8(value%24), value/24
	day := value
	if day != 0 {
		return fmt.Sprintf("%s%dd%02d:%02d:%02d.%09d", sign, day, hour, minute, second, nano)
	}
	if hour != 0 {
		return fmt.Sprintf("%s%d:%02d:%02d.%09d", sign, hour, minute, second, nano)
	}
	return fmt.Sprintf("%s%d:%02d.%09d", sign, minute, second, nano)
}

// Duration formats the rounded duration
func Duration(d time.Duration) string {
	value := int64(d / time.Nanosecond)
	sign := ""
	s := int64(1)
	if value < 0 {
		sign = "-"
		s = int64(-1)
	}
	var nano, micro, milli uint16
	var second, minute, hour uint8
	nano, value = uint16(value%1000*s), value/1000*s
	micro, value = uint16(value%1000), value/1000
	milli, value = uint16(value%1000), value/1000
	second, value = uint8(value%60), value/60
	minute, value = uint8(value%60), value/60
	hour, value = uint8(value%24), value/24
	day := value
	if day != 0 {
		seconds := float64(second) + float64(milli)*1e-3 + float64(micro)*1e-6 + float64(nano)*1e-9
		return fmt.Sprintf("%s%dd%02d:%02d:%02.0f", sign, day, hour, minute, seconds)
	}
	if hour != 0 {
		seconds := float64(second) + float64(milli)*1e-3 + float64(micro)*1e-6 + float64(nano)*1e-9
		return fmt.Sprintf("%s%d:%02d:%02.0f", sign, hour, minute, seconds)
	}
	if minute != 0 {
		seconds := float64(second) + float64(milli)*1e-3 + float64(micro)*1e-6 + float64(nano)*1e-9
		if minute >= 10 {
			return fmt.Sprintf("%s%d:%02.0f", sign, minute, seconds)
		}
		return fmt.Sprintf("%s%d:%04.1f", sign, minute, seconds)
	}
	if second != 0 {
		if second >= 10 {
			seconds := float64(second) + float64(milli)*1e-3 + float64(micro)*1e-6 + float64(nano)*1e-9
			return fmt.Sprintf("%s%.2fs", sign, seconds)
		}
		millis := float64(second)*1e+3 + float64(milli) + float64(micro)*1e-3 + float64(nano)*1e-6
		return fmt.Sprintf("%s%.0fms", sign, millis)
	}
	if milli != 0 {
		millis := float64(milli) + float64(micro)*1e-3 + float64(nano)*1e-6
		if milli >= 100 {
			return fmt.Sprintf("%s%.1fms", sign, millis)
		}
		if milli >= 10 {
			return fmt.Sprintf("%s%.2fms", sign, millis)
		}
		return fmt.Sprintf("%s%.3fms", sign, millis)
	}
	if micro != 0 {
		if micro >= 100 {
			micros := float64(micro) + float64(nano)*1e-3
			return fmt.Sprintf("%s%.1fµs", sign, micros)
		}
		if micro >= 10 {
			micros := float64(micro) + float64(nano)*1e-3
			return fmt.Sprintf("%s%.2fµs", sign, micros)
		}
		nanos := float64(micro)*1e+3 + float64(nano)
		return fmt.Sprintf("%s%.0fns", sign, nanos)
	}
	if nano != 0 {
		return fmt.Sprintf("%s%dns", sign, nano)
	}
	return fmt.Sprintf("%s0", sign)
}
