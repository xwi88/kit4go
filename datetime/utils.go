package datetime

import (
	"fmt"
	"strings"
	"time"
)

// GetDateTimeLayoutWithFormatAndZoneOffset default 2006-01-02T15:04:05%v:00
func GetDateTimeLayoutWithFormatAndZoneOffset(format string, zoneOffset int) string {
	var zone string
	if zoneOffset <= -24 || zoneOffset >= 24 {
		zoneOffset = 0
	}
	if zoneOffset < 0 {
		zone = fmt.Sprintf("-%02d", -zoneOffset)
	} else {
		zone = fmt.Sprintf("+%02d", zoneOffset)
	}
	if len(format) == 0 {
		format = FormatLayoutDateTimeISO8601WithZone
	}
	if !strings.Contains(format, "%") {
		return format
	}
	return fmt.Sprintf(format, zone)
}

// GetDateTimeLayoutISO8601WithZoneOffset 2006-01-02T15:04:05%v:00
func GetDateTimeLayoutISO8601WithZoneOffset(zoneOffset int) string {
	return GetDateTimeLayoutWithFormatAndZoneOffset(FormatLayoutDateTimeISO8601WithZone, zoneOffset)
}

// GetDateTimeLayoutISO8601WithFormatAndZoneOffsetZoneMid 2006-01-02T15:04:05%v00
func GetDateTimeLayoutISO8601WithFormatAndZoneOffsetZoneMid(zoneOffset int) string {
	return GetDateTimeLayoutWithFormatAndZoneOffset(FormatLayoutDateTimeISO8601WithZoneMid, zoneOffset)
}

// NowUnix return now Unix time, second
func NowUnix() int64 {
	return time.Now().Unix()
}

// NowUnixMillisecond return now Millisecond time, millisecond
func NowUnixMillisecond() int64 {
	return int64(time.Now().Nanosecond() / 1000000)
}

// NowUnixMicrosecond return now Microsecond time, microsecond
func NowUnixMicrosecond() int64 {
	return int64(time.Now().Nanosecond() / 1000)
}

// NowUnixNano return now UnixNano time, unix nano second
func NowUnixNano() int64 {
	return time.Now().UnixNano()
}

// NowTimestampStr with layout
func NowTimestampStr(layout string) string {
	return time.Now().Format(layout)
}

// DeltaHours return the hours for t2 - t1
func DeltaHours(t1, t2 time.Time) float64 {
	return t2.Sub(t1).Hours() / 24.0
}

// DeltaDays return the real days between t1 and t2
// t2 - t1
func DeltaDays(t1, t2 time.Time) int {
	if t1.Equal(t2) {
		return 0
	}
	t1ds := t1.Format(LayoutDateISO8601)
	t2ds := t2.Format(LayoutDateISO8601)
	t1dt, _ := ParseTime(LayoutDateISO8601, t1ds)
	t2dt, _ := ParseTime(LayoutDateISO8601, t2ds)
	return int(t2dt.Sub(t1dt).Hours()/24.0) + 1
}

// AddDay add day for the special time
func AddDay(day int, d time.Time) time.Time {
	return d.AddDate(0, 0, day)
}

func time2Str(layout string, d time.Time) string {
	if len(layout) == 0 {
		layout = LayoutDateTime
	}
	return d.Format(layout)
}

// GetTimeStr get time str for the special time
func GetTimeStr(layout string, d time.Time) string {
	return time2Str(layout, d)
}

// ParseTime parse time with local timezone
func ParseTime(layout, value string) (time.Time, error) {
	return ParseTimeWithLocation(layout, value, time.Local)
}

// ParseTimeWithLocation ...
func ParseTimeWithLocation(layout, value string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.Local
	}
	return time.ParseInLocation(layout, value, loc)
}

// GetFirstDateOfMonth get the first datetime for this month
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetStartTime(d)
}

// GetLastDateOfMonth get the last datetime for this month
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// GetStartTime get the start time for the special day, ex: 2006-06-01T00:00:00.999+08:00
func GetStartTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// GetStartTimeStr get start time str for the special time
func GetStartTimeStr(layout string, d time.Time) string {
	if len(layout) == 0 {
		layout = LayoutDateTime
	}
	return GetStartTime(d).Format(layout)
}

// GetEndTime get the end time for the special day, ex: 2019-06-01T23:59:59.999+08:00
func GetEndTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 1e9-1, d.Location())
}

// GetEndTimeStr get end time str for the special time
func GetEndTimeStr(layout string, d time.Time) string {
	if len(layout) == 0 {
		layout = LayoutDateTime
	}
	return GetEndTime(d).Format(layout)
}

// GetStartEndTimeStr get start and end time str for the special day
func GetStartEndTimeStr(layout string, d time.Time) (start, end string) {
	if len(layout) == 0 {
		layout = LayoutDateTime
	}
	start = GetStartTime(d).Format(layout)
	end = GetEndTime(d).Format(layout)
	return
}

// GetStartEndTimeStrWithZone get start and end time str with zone for the special day
func GetStartEndTimeStrWithZone(layout string, d time.Time) (start, end string) {
	if len(layout) == 0 {
		layout = LayoutDateTimeISO8601Zone
	}
	start = GetStartTime(d).Format(layout)
	end = GetEndTime(d).Format(layout)
	return
}

// GetNowWithZone get now time with special location
// time format: ISO8601:2004 2004-05-03T17:30:08+08:00
// go format: 2006-01-02T15:04:05+00:00
func GetNowWithZone(loc *time.Location) time.Time {
	now := time.Now()
	if loc == nil {
		loc = time.Local
	}
	if t, err := time.ParseInLocation(LayoutDateTime, time2Str(LayoutDateTime, now), loc); err == nil {
		return t
	}
	return now
}

// GetNowWithZoneAndLayout get now time with special location and special layout
func GetNowWithZoneAndLayout(loc *time.Location, layout string) time.Time {
	now := time.Now()
	if loc == nil {
		loc = time.Local
	}
	if layout == "" {
		layout = LayoutDateTime
	}
	if t, err := time.ParseInLocation(layout, time2Str(layout, now), loc); err == nil {
		return t
	}
	return now
}
