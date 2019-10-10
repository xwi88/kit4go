package datetime

import (
	"testing"
)

func TestGetDateTimeLayoutISO8601WithFormatAndZoneOffset(t *testing.T) {
	outputFormat := ""
	formats := []string{LayoutTime, LayoutTimeShort,LayoutDateISO8601, LayoutDateISO8601Short,
		LayoutDateTime, LayoutDateTimeShort}
	zoneOffsets := []int{0, +8}
	for i := range formats {
		for j := range zoneOffsets {
			offset := zoneOffsets[j]
			outputFormat = GetDateTimeLayoutISO8601WithFormatAndZoneOffset(formats[i], offset)
			t.Logf("format:%v, zoneOffset:%v, output:%s", formats[i], offset, outputFormat)
		}
	}

	formats = []string{
		FormatLayoutTimeISO8601WithZone,
		FormatLayoutTimeISO8601WithZoneMid,
		FormatLayoutTimeISO8601WithZoneShort,
		FormatLayoutDateTimeISO8601WithZone,
		FormatLayoutDateTimeISO8601WithZoneMid,
		FormatLayoutDateTimeISO8601WithZoneShort,
		FormatLayoutDateTimeISO8601ShortWithZone,
		FormatLayoutDateTimeISO8601ShortWithZoneMid,
		FormatLayoutDateTimeISO8601ShortWithZoneShort,
	}
	zoneOffsets = []int{0, +8, -8, +24}
	for i := range formats {
		for j := range zoneOffsets {
			offset := zoneOffsets[j]
			outputFormat = GetDateTimeLayoutISO8601WithFormatAndZoneOffset(formats[i], offset)
			t.Logf("format:%v, zoneOffset:%v, output:%s", formats[i], offset, outputFormat)
		}
	}
}
