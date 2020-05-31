package datetime

// layout format with ISO8601:2004
const (
	// common use layout format without zone
	LayoutTime             = "15:04:05"            // 8bit
	LayoutTimeShort        = "150405"              // 6bit
	LayoutDateISO8601      = "2006-01-02"          // 10bit
	LayoutDateISO8601Short = "20060102"            //  8bit
	LayoutDateTime         = "2006-01-02 15:04:05" // 19bit
	LayoutDateTimeShort    = "20060102150405"      // 14bit

	DefaultLayoutTime     = LayoutTime
	DefaultLayoutDate     = LayoutDateISO8601
	DefaultLayoutDateTime = LayoutDateTime

	DefaultZoneOffset = "+00"

	// common use layout format with zone
	FormatLayoutTimeISO8601WithZone               = "15:04:05%v:00"            // 14bit zone +xx:00
	FormatLayoutTimeISO8601WithZoneMid            = "150405%v00"               // 11bit zone +xx00
	FormatLayoutTimeISO8601WithZoneShort          = "150405%v"                 //  9bit zone +xx
	FormatLayoutDateTimeISO8601WithZone           = "2006-01-02T15:04:05%v:00" // 25bit zone +xx:00
	FormatLayoutDateTimeISO8601WithZoneMid        = "2006-01-02T15:04:05%v00"  // 24bit zone +xx:00
	FormatLayoutDateTimeISO8601WithZoneShort      = "2006-01-02T15:04:05%v"    // 22bit zone +xx:00
	FormatLayoutDateTimeISO8601ShortWithZone      = "20060102T150405%v:00"     // 21bit zone +xx00
	FormatLayoutDateTimeISO8601ShortWithZoneMid   = "20060102T150405%v00"      // 20bit zone +xx00
	FormatLayoutDateTimeISO8601ShortWithZoneShort = "20060102T150405%v"        // 18bit zone +xx

	// zone +00:00 | +0000 | +00
	LayoutTimeISO8601                   = "15:04:05+00:00"            // 14bit
	LayoutTimeISO8601Mid                = "150405+0000"               // 11bit
	LayoutTimeISO8601Short              = "150405+00"                 //  9bit
	LayoutDateTimeISO8601Zone           = "2006-01-02T15:04:05+00:00" // 25bit zone +00:00
	LayoutDateTimeISO8601ZoneMid        = "2006-01-02T15:04:05+0000"  // 24bit zone +0000
	LayoutDateTimeISO8601ZoneShort      = "2006-01-02T15:04:05+00"    // 25bit zone +00
	LayoutDateTimeISO8601ShortZone      = "20060102T150405+00:00"     // 21bit zone +0000
	LayoutDateTimeISO8601ShortZoneMid   = "20060102T150405+0000"      // 20bit zone +0000
	LayoutDateTimeISO8601ShortZoneShort = "20060102T150405+00"        // 18bit zone +00
)

// common layout format with special zone 08
const (
	// zone +08:00 | +0800 | +08
	// China time and datetime layout format with ISO8601:2004
	LayoutTimeISO8601ZoneP8               = "15:04:05+08:00"            // 14bit zone +08:00
	LayoutTimeISO8601ZoneP8Mid            = "150405+0800"               // 11bit zone +0800
	LayoutTimeISO8601ZoneP8Short          = "150405+08"                 //  9bit zone +08
	LayoutDateTimeISO8601ZoneP8           = "2006-01-02T15:04:05+08:00" // 25bit zone +08:00
	LayoutDateTimeISO8601ZoneP8Mid        = "2006-01-02T15:04:05+0800"  // 24bit zone +08:00
	LayoutDateTimeISO8601ZoneP8Short      = "2006-01-02T15:04:05+08"    // 22bit zone +08:00
	LayoutDateTimeISO8601ShortZoneP8      = "20060102T150405+08:00"     // 21bit zone +0800
	LayoutDateTimeISO8601ShortZoneP8Mid   = "20060102T150405+0800"      // 20bit zone +0800
	LayoutDateTimeISO8601ShortZoneP8Short = "20060102T150405+08"        // 18bit zone +08
)
