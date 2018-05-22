package gopifinder

import (
	"strconv"
	"strings"
	"time"
)

// ConvToInt64 converts the string value to a Int64 value.
func ConvToInt64(s string, def int64) int64 {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return def
	}
	return n
}

// ConvToInt converts the string value to a Int value.
func ConvToInt(s string, def int) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return def
	}
	return n
}

// ConvToDate converts the WMIC date format to a standard date
func ConvToDate(s string) time.Time {
	s = strings.TrimSpace(s)
	y := ConvToInt(s[0:4], 0)
	m := time.Month(ConvToInt(s[4:6], 0))
	d := ConvToInt(s[6:8], 0)
	h := ConvToInt(s[8:10], 0)
	n := ConvToInt(s[10:12], 0)
	c := ConvToInt(s[12:14], 0)

	return time.Date(y, m, d, h, n, c, 0, time.Local)
}
