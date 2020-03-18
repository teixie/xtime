package xtime

import (
	"fmt"
	"regexp"
	"time"
)

const (
	LayoutYmdHis = "2006-01-02 15:04:05"
	LayoutYmd    = "2006-01-02"
)

var (
	local *time.Location
)

// 设置时区
func SetLocation(loc *time.Location) {
	local = loc
}

// 获得时区
func GetLocation() *time.Location {
	if local != nil {
		return local
	}
	return time.Local
}

// 获得格式化时间：2006-01-02 15:04:05
func YmdHis(t time.Time) string {
	return t.Format(LayoutYmdHis)
}

// 获得格式化时间：2006-01-02
func Ymd(t time.Time) string {
	return t.Format(LayoutYmd)
}

// 时间格式化，支持"Y-m-d H:i:s"、"YYYY-mm-dd HH:ii:ss"等形式，当不包含任何"YymdHis"字符时将使用原生格式化方法
func Format(t time.Time, fmtStr string) string {
	exists, err := regexp.Match("[YymdHis]+", []byte(fmtStr))
	if err == nil && !exists {
		return t.Format(fmtStr)
	}

	timeStr := t.String()
	o := map[string]string{
		"Y+": timeStr[0:4],
		"y+": timeStr[2:4],
		"m+": timeStr[5:7],
		"d+": timeStr[8:10],
		"H+": timeStr[11:13],
		"i+": timeStr[14:16],
		"s+": timeStr[17:19],
	}
	for k, v := range o {
		re, _ := regexp.Compile(k)
		fmtStr = re.ReplaceAllString(fmtStr, v)
	}
	return fmtStr
}

// 获得当前时间
func Now() time.Time {
	return time.Now().In(GetLocation())
}

// 获得今天的开始时间
func Today() time.Time {
	t := time.Now().In(GetLocation())
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetLocation())
}

// 获得明天的开始时间
func Tomorrow(args ...time.Time) time.Time {
	t := findOrNow(args...).Add(24 * time.Hour)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetLocation())
}

// 获得昨天的开始时间
func Yesterday(args ...time.Time) time.Time {
	t := findOrNow(args...).Add(-24 * time.Hour)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetLocation())
}

// 所在天的开始时间，如："2006-01-01 00:00:00"
func StartOfDay(args ...time.Time) time.Time {
	t := findOrNow(args...)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetLocation())
}

// 所在天的结束时间，如："2006-01-01 23:59:59"
func EndOfDay(args ...time.Time) time.Time {
	t := findOrNow(args...)
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, GetLocation())
}

// 当前时间所在星期的开始时间，例："2006-01-02 00:00:00"
func StartOfWeek(args ...time.Time) time.Time {
	now := findOrNow(args...)
	t := now.Add(-(time.Duration(now.Weekday()) - 1) * 24 * time.Hour)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetLocation())
}

// 当前时间所在星期的结束时间，例："2006-01-02 23:59:59"
func EndOfWeek(args ...time.Time) time.Time {
	now := findOrNow(args...)
	t := now.Add((7 - time.Duration(now.Weekday())) * 24 * time.Hour)
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, GetLocation())
}

// 当前时间所在月的开始时间，例："2016-01-01 00:00:00"
func StartOfMonth(args ...time.Time) time.Time {
	now := findOrNow(args...)
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, GetLocation())
}

// 当前时间所在月的结束时间，例："2016-01-31 23:59:59"
func EndOfMonth(args ...time.Time) time.Time {
	now := findOrNow(args...)
	next := StartOfMonth(now).Add(31 * 24 * time.Hour)
	return time.Unix(StartOfMonth(next).Unix()-1, 0)
}

// 当前时间所在年的开始时间，例："2016-01-01 00:00:00"
func StartOfYear(args ...time.Time) time.Time {
	now := findOrNow(args...)
	return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, GetLocation())
}

// 当前时间所在年的结束时间，例："2016-12-31 23:59:59"
func EndOfYear(args ...time.Time) time.Time {
	now := findOrNow(args...)
	return time.Date(now.Year(), 12, 31, 23, 59, 59, 0, GetLocation())
}

// 解析时间，支持time.Time/时间字符串/时间戳
func Parse(t interface{}) (time.Time, error) {
	if t == nil {
		return time.Time{}, nil
	}

	switch t.(type) {
	case time.Time:
		return t.(time.Time), nil
	case string:
		return time.ParseInLocation(LayoutYmdHis, t.(string), GetLocation())
	case int:
		return time.Unix(int64(t.(int)), 0), nil
	case int64:
		return time.Unix(t.(int64), 0), nil
	}

	return time.Time{}, fmt.Errorf("type not support")
}

// 获取第一个或者当前时间
func findOrNow(args ...time.Time) time.Time {
	if len(args) > 0 {
		return args[0]
	}
	return Now()
}
