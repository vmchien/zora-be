package format

import (
	"time"

	"vn.vato.zora.be.api/pkg/constant"
)

func TimeToString(t time.Time, formats ...string) string {
	if len(formats) == 0 {
		return t.Format(constant.DEFAULT_TIME_FORMAT)
	}
	return t.Format(formats[0])
}

func StringToTime(timeStr string, formats ...string) (time.Time, error) {
	if len(formats) == 0 {
		return time.Parse(constant.DEFAULT_TIME_FORMAT, timeStr)
	}
	return time.Parse(formats[0], timeStr)
}

func StringToTimeWithTZ(timeStr string, tz *time.Location, formats ...string) (time.Time, error) {
	layout := constant.DEFAULT_TIME_FORMAT
	if len(formats) > 0 {
		layout = formats[0]
	}
	return time.ParseInLocation(layout, timeStr, tz)
}

func TryParseStringToTimeWithTZ(timeStr string, tz *time.Location, formats ...string) (time.Time, error) {
	rs, err := StringToTimeWithTZ(timeStr, tz, formats...)
	if err != nil {
		return StringToTimeWithTZ(timeStr, tz, constant.DEFAULT_TIME_FORMAT)
	}
	return rs, nil
}

//
// func TimeToVNString(t time.Time, formats ...string) string {
//
// 	if len(formats) == 0 {
// 		return t.Format(constant.DEFAULT_TIME_FORMAT)
// 	}
// 	return t.Format(formats[0])
// }
