// gosnippets (c) 2023-2025 He Xian <hexian000@outlook.com>
// This code is licensed under MIT license (see LICENSE for details)

package formats_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/hexian000/gosnippets/formats"
)

func TestSIPrefix(t *testing.T) {
	zero := 0.0
	m := 2.0 / 3.0
	cases := [...]struct {
		in float64
		s  string
	}{
		{math.NaN(), "nan"}, {-math.NaN(), "-nan"},
		{math.Inf(1), "inf"}, {math.Inf(-1), "-inf"},
		{zero, "0"}, {-zero, "-0"},
		{1e-39 * m, "6.67e-40"}, {1e-38 * m, "6.67e-39"},
		{1e-37 * m, "6.67e-38"}, {1e-36 * m, "6.67e-37"},
		{1e-35 * m, "6.67e-36"}, {1e-34 * m, "6.67e-35"},
		{1e-33 * m, "6.67e-34"}, {1e-32 * m, "6.67e-33"},
		{1e-31 * m, "6.67e-32"}, {1e-30 * m, "6.67e-31"},
		{1e-29 * m, "6.67q"}, {1e-28 * m, "66.7q"},
		{1e-27 * m, "667q"}, {1e-26 * m, "6.67r"},
		{1e-25 * m, "66.7r"}, {1e-24 * m, "667r"},
		{1e-23 * m, "6.67y"}, {1e-22 * m, "66.7y"},
		{1e-21 * m, "667y"}, {1e-20 * m, "6.67z"},
		{1e-19 * m, "66.7z"}, {1e-18 * m, "667z"},
		{1e-17 * m, "6.67a"}, {1e-16 * m, "66.7a"},
		{1e-15 * m, "667a"}, {1e-14 * m, "6.67f"},
		{1e-13 * m, "66.7f"}, {1e-12 * m, "667f"},
		{1e-11 * m, "6.67p"}, {1e-10 * m, "66.7p"},
		{1e-09 * m, "667p"}, {1e-08 * m, "6.67n"},
		{1e-07 * m, "66.7n"}, {1e-06 * m, "667n"},
		{1e-05 * m, "6.67μ"}, {1e-04 * m, "66.7μ"},
		{1e-03 * m, "667μ"}, {1e-02 * m, "6.67m"},
		{1e-01 * m, "66.7m"}, {1e+00 * m, "667m"},
		{1e+01 * m, "6.67"}, {1e+02 * m, "66.7"},
		{1e+03 * m, "667"}, {1e+04 * m, "6.67k"},
		{1e+05 * m, "66.7k"}, {1e+06 * m, "667k"},
		{1e+07 * m, "6.67M"}, {1e+08 * m, "66.7M"},
		{1e+09 * m, "667M"}, {1e+10 * m, "6.67G"},
		{1e+11 * m, "66.7G"}, {1e+12 * m, "667G"},
		{1e+13 * m, "6.67T"}, {1e+14 * m, "66.7T"},
		{1e+15 * m, "667T"}, {1e+16 * m, "6.67P"},
		{1e+17 * m, "66.7P"}, {1e+18 * m, "667P"},
		{1e+19 * m, "6.67E"}, {1e+20 * m, "66.7E"},
		{1e+21 * m, "667E"}, {1e+22 * m, "6.67Z"},
		{1e+23 * m, "66.7Z"}, {1e+24 * m, "667Z"},
		{1e+25 * m, "6.67Y"}, {1e+26 * m, "66.7Y"},
		{1e+27 * m, "667Y"}, {1e+28 * m, "6.67R"},
		{1e+29 * m, "66.7R"}, {1e+30 * m, "667R"},
		{1e+31 * m, "6.67Q"}, {1e+32 * m, "6.67e+31"},
		{1e+33 * m, "6.67e+32"}, {1e+34 * m, "6.67e+33"},
		{1e+35 * m, "6.67e+34"}, {1e+36 * m, "6.67e+35"},
		{1e+37 * m, "6.67e+36"}, {1e+38 * m, "6.67e+37"},
	}
	for _, c := range cases {
		value := c.in
		result := formats.SIPrefix(value)
		if result != c.s {
			t.Fatalf("expect \"%s\", got \"%s\"", c.s, result)
		}
		fmt.Printf("|%16e|%16s|\n", value, result)
		// fmt.Printf("{%.0e*m, \"%s\"},\n", value/m, result)
	}
}

func TestIECBytes(t *testing.T) {
	zero := 0.0
	cases := [...]struct {
		in float64
		s  string
	}{
		{math.NaN(), "nan"}, {-math.NaN(), "-nan"},
		{math.Inf(1), "inf"}, {math.Inf(-1), "-inf"},
		{zero, "0"}, {-zero, "-0"},
		{1e+00, "1B"}, {1e+01, "10B"}, {1e+02, "100B"}, {1e+03, "1000B"},
		{2e+03, "2000B"}, {4e+03, "3.91KiB"}, {6e+03, "5.86KiB"}, {8e+03, "7.81KiB"},
		{1e+04, "9.77KiB"}, {1e+05, "97.7KiB"}, {1e+06, "977KiB"}, {1e+07, "9.54MiB"},
		{1e+08, "95.4MiB"}, {1e+09, "954MiB"}, {1e+10, "9.31GiB"}, {1e+11, "93.1GiB"},
		{1e+12, "931GiB"}, {1e+13, "9.09TiB"}, {1e+14, "90.9TiB"}, {1e+15, "909TiB"},
		{1e+16, "8.88PiB"}, {1e+17, "88.8PiB"}, {1e+18, "888PiB"}, {1e+19, "8.67EiB"},
		{1e+20, "86.7EiB"}, {1e+21, "867EiB"}, {1e+22, "8.47ZiB"}, {1e+23, "84.7ZiB"},
		{1e+24, "847ZiB"}, {1e+25, "8.27YiB"}, {1e+26, "82.7YiB"}, {1e+27, "827YiB"},
		{1e+28, "8272YiB"}, {1e+29, "82718YiB"},
	}
	for _, c := range cases {
		value := c.in
		result := formats.IECBytes(value)
		if result != c.s {
			t.Fatalf("expect \"%s\", got \"%s\"", c.s, result)
		}
		fmt.Printf("|%16e|%16s|\n", value, result)
		// fmt.Printf("{%.0e, \"%s\"},\n", value, result)
	}
}

func TestDuration(t *testing.T) {
	cases := [...]struct {
		in                         time.Duration
		secs, millis, nanos, human string
	}{
		{1 * time.Minute, "1:00", "1:00.000", "1:00.000000000", "1:00.0"},
		{math.MaxInt64, "106751d23:47:16", "106751d23:47:16.854", "106751d23:47:16.854775807", "106751d23:47:17"},
		{math.MinInt64, "-106751d23:47:16", "-106751d23:47:16.854", "-106751d23:47:16.854775808", "-106751d23:47:17"},
		{0, "0:00", "0:00.000", "0:00.000000000", "0"},
		{1, "0:00", "0:00.000", "0:00.000000001", "1ns"},
		{-1, "-0:00", "-0:00.000", "-0:00.000000001", "-1ns"},
		{4, "0:00", "0:00.000", "0:00.000000004", "4ns"},
		{-4, "-0:00", "-0:00.000", "-0:00.000000004", "-4ns"},
		{16, "0:00", "0:00.000", "0:00.000000016", "16ns"},
		{-16, "-0:00", "-0:00.000", "-0:00.000000016", "-16ns"},
		{64, "0:00", "0:00.000", "0:00.000000064", "64ns"},
		{-64, "-0:00", "-0:00.000", "-0:00.000000064", "-64ns"},
		{256, "0:00", "0:00.000", "0:00.000000256", "256ns"},
		{-256, "-0:00", "-0:00.000", "-0:00.000000256", "-256ns"},
		{1024, "0:00", "0:00.000", "0:00.000001024", "1024ns"},
		{-1024, "-0:00", "-0:00.000", "-0:00.000001024", "-1024ns"},
		{4096, "0:00", "0:00.000", "0:00.000004096", "4096ns"},
		{-4096, "-0:00", "-0:00.000", "-0:00.000004096", "-4096ns"},
		{16384, "0:00", "0:00.000", "0:00.000016384", "16.38µs"},
		{-16384, "-0:00", "-0:00.000", "-0:00.000016384", "-16.38µs"},
		{65536, "0:00", "0:00.000", "0:00.000065536", "65.54µs"},
		{-65536, "-0:00", "-0:00.000", "-0:00.000065536", "-65.54µs"},
		{262144, "0:00", "0:00.000", "0:00.000262144", "262.1µs"},
		{-262144, "-0:00", "-0:00.000", "-0:00.000262144", "-262.1µs"},
		{1048576, "0:00", "0:00.001", "0:00.001048576", "1.049ms"},
		{-1048576, "-0:00", "-0:00.001", "-0:00.001048576", "-1.049ms"},
		{4194304, "0:00", "0:00.004", "0:00.004194304", "4.194ms"},
		{-4194304, "-0:00", "-0:00.004", "-0:00.004194304", "-4.194ms"},
		{16777216, "0:00", "0:00.016", "0:00.016777216", "16.78ms"},
		{-16777216, "-0:00", "-0:00.016", "-0:00.016777216", "-16.78ms"},
		{67108864, "0:00", "0:00.067", "0:00.067108864", "67.11ms"},
		{-67108864, "-0:00", "-0:00.067", "-0:00.067108864", "-67.11ms"},
		{268435456, "0:00", "0:00.268", "0:00.268435456", "268.4ms"},
		{-268435456, "-0:00", "-0:00.268", "-0:00.268435456", "-268.4ms"},
		{1073741824, "0:01", "0:01.073", "0:01.073741824", "1074ms"},
		{-1073741824, "-0:01", "-0:01.073", "-0:01.073741824", "-1074ms"},
		{4294967296, "0:04", "0:04.294", "0:04.294967296", "4295ms"},
		{-4294967296, "-0:04", "-0:04.294", "-0:04.294967296", "-4295ms"},
		{17179869184, "0:17", "0:17.179", "0:17.179869184", "17.18s"},
		{-17179869184, "-0:17", "-0:17.179", "-0:17.179869184", "-17.18s"},
		{68719476736, "1:08", "1:08.719", "1:08.719476736", "1:08.7"},
		{-68719476736, "-1:08", "-1:08.719", "-1:08.719476736", "-1:08.7"},
		{274877906944, "4:34", "4:34.877", "4:34.877906944", "4:34.9"},
		{-274877906944, "-4:34", "-4:34.877", "-4:34.877906944", "-4:34.9"},
		{1099511627776, "18:19", "18:19.511", "18:19.511627776", "18:20"},
		{-1099511627776, "-18:19", "-18:19.511", "-18:19.511627776", "-18:20"},
		{4398046511104, "1:13:18", "1:13:18.046", "1:13:18.046511104", "1:13:18"},
		{-4398046511104, "-1:13:18", "-1:13:18.046", "-1:13:18.046511104", "-1:13:18"},
		{17592186044416, "4:53:12", "4:53:12.186", "4:53:12.186044416", "4:53:12"},
		{-17592186044416, "-4:53:12", "-4:53:12.186", "-4:53:12.186044416", "-4:53:12"},
		{70368744177664, "19:32:48", "19:32:48.744", "19:32:48.744177664", "19:32:49"},
		{-70368744177664, "-19:32:48", "-19:32:48.744", "-19:32:48.744177664", "-19:32:49"},
		{281474976710656, "3d06:11:14", "3d06:11:14.976", "3d06:11:14.976710656", "3d06:11:15"},
		{-281474976710656, "-3d06:11:14", "-3d06:11:14.976", "-3d06:11:14.976710656", "-3d06:11:15"},
		{1125899906842624, "13d00:44:59", "13d00:44:59.906", "13d00:44:59.906842624", "13d00:44:60"},
		{-1125899906842624, "-13d00:44:59", "-13d00:44:59.906", "-13d00:44:59.906842624", "-13d00:44:60"},
		{4503599627370496, "52d02:59:59", "52d02:59:59.627", "52d02:59:59.627370496", "52d02:59:60"},
		{-4503599627370496, "-52d02:59:59", "-52d02:59:59.627", "-52d02:59:59.627370496", "-52d02:59:60"},
		{18014398509481984, "208d11:59:58", "208d11:59:58.509", "208d11:59:58.509481984", "208d11:59:59"},
		{-18014398509481984, "-208d11:59:58", "-208d11:59:58.509", "-208d11:59:58.509481984", "-208d11:59:59"},
		{72057594037927936, "833d23:59:54", "833d23:59:54.037", "833d23:59:54.037927936", "833d23:59:54"},
		{-72057594037927936, "-833d23:59:54", "-833d23:59:54.037", "-833d23:59:54.037927936", "-833d23:59:54"},
		{288230376151711744, "3335d23:59:36", "3335d23:59:36.151", "3335d23:59:36.151711744", "3335d23:59:36"},
		{-288230376151711744, "-3335d23:59:36", "-3335d23:59:36.151", "-3335d23:59:36.151711744", "-3335d23:59:36"},
		{1152921504606846976, "13343d23:58:24", "13343d23:58:24.606", "13343d23:58:24.606846976", "13343d23:58:25"},
		{-1152921504606846976, "-13343d23:58:24", "-13343d23:58:24.606", "-13343d23:58:24.606846976", "-13343d23:58:25"},
	}
	for _, c := range cases {
		secs := formats.DurationSeconds(c.in)
		millis := formats.DurationMillis(c.in)
		nanos := formats.DurationNanos(c.in)
		human := formats.Duration(c.in)
		if secs != c.secs {
			t.Fatalf("expect \"%s\", got \"%s\"", c.secs, secs)
		}
		if millis != c.millis {
			t.Fatalf("expect \"%s\", got \"%s\"", c.millis, millis)
		}
		if nanos != c.nanos {
			t.Fatalf("expect \"%s\", got \"%s\"", c.nanos, nanos)
		}
		if human != c.human {
			t.Fatalf("expect \"%s\", got \"%s\"", c.human, human)
		}
		fmt.Printf("|%16s|%20s|%26s|%16s|\n", secs, millis, nanos, human)
		// fmt.Printf("{%d, \"%s\", \"%s\", \"%s\", \"%s\"},\n", c.in,
		// 	secs, millis, nanos, human)
	}
}
