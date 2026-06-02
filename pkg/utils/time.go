// Package utils provides timezone-safe time helpers.
//
// Design principles:
//   - Never use time.Local (it is environment-dependent and often UTC in containers).
//   - Always convert to a target location first, then compute day boundaries.
//   - Prefer end-exclusive ranges [start, end) for Elasticsearch (gte/lt) to avoid precision issues.
package utils

import "time"

// VNLocation returns the Asia/Ho_Chi_Minh location.
//
// Note: time.LoadLocation depends on tzdata being available in the runtime image.
// If your container is scratch/distroless without tzdata, load may fail.
// In that case, use a fixed zone: time.FixedZone("Asia/Ho_Chi_Minh", 7*3600).
func VNLocation() (*time.Location, error) {
	return time.LoadLocation("Asia/Ho_Chi_Minh")
}

// MustVNLocation is a convenience for initialization paths where failing fast is acceptable.
func MustVNLocation() *time.Location {
	loc, err := VNLocation()
	if err != nil {
		// Fallback to fixed zone to avoid crashing in minimal images.
		// This is safe for Ho Chi Minh because it does not observe DST.
		return time.FixedZone("Asia/Ho_Chi_Minh", 7*3600)
	}
	return loc
}

// ToLocation converts a time instant to the given location.
func ToLocation(t time.Time, loc *time.Location) time.Time {
	if loc == nil {
		// Caller mistake; keep original location rather than guessing.
		return t
	}
	return t.In(loc)
}

// StartOfDay returns 00:00:00.000 of the day that contains t, evaluated in loc.
func StartOfDay(t time.Time, loc *time.Location) time.Time {
	tv := ToLocation(t, loc)
	return time.Date(tv.Year(), tv.Month(), tv.Day(), 0, 0, 0, 0, loc)
}

// EndOfDayExclusive returns the start of the next day (00:00:00.000), evaluated in loc.
// This is the recommended "end" bound for range queries using lt.
func EndOfDayExclusive(t time.Time, loc *time.Location) time.Time {
	tv := ToLocation(t, loc)
	return time.Date(tv.Year(), tv.Month(), tv.Day()+1, 0, 0, 0, 0, loc)
}

// EndOfDayInclusive returns 23:59:59.999999999 of the day that contains t, evaluated in loc.
// Prefer EndOfDayExclusive for Elasticsearch ranges (gte/lt) to avoid precision pitfalls.
func EndOfDayInclusive(t time.Time, loc *time.Location) time.Time {
	return EndOfDayExclusive(t, loc).Add(-time.Nanosecond)
}

// TimeOfDay returns t represented in loc, truncating nanoseconds to zero.
// This is useful when you want stable string serialization at second precision.
func TimeOfDay(t time.Time, loc *time.Location) time.Time {
	tv := ToLocation(t, loc)
	return time.Date(tv.Year(), tv.Month(), tv.Day(), tv.Hour(), tv.Minute(), tv.Second(), 0, loc)
}

// DateTimeRange returns [start, end) for the day that contains t, evaluated in loc.
// Intended for Elasticsearch:
//   - gte = start
//   - lt  = end
func DateTimeRange(t time.Time, loc *time.Location) (start time.Time, endExclusive time.Time) {
	start = TimeOfDay(t, loc)
	endExclusive = EndOfDayExclusive(t, loc)
	return
}

// DayRange returns [start, end) for the day that contains t, evaluated in loc.
// Intended for Elasticsearch:
//   - gte = start
//   - lt  = end
func DayRange(t time.Time, loc *time.Location) (start time.Time, endExclusive time.Time) {
	start = StartOfDay(t, loc)
	endExclusive = EndOfDayExclusive(t, loc)
	return
}

// DayRangeUTC returns [startUTC, endUTC) for the day that contains t when interpreted in loc.
// This is typically what you want if your ES date fields are stored/queried in UTC ("Z").
// Example:
//
//	startUTC, endUTC := DayRangeUTC(targetDate, vnLoc)
//	range: gte=startUTC.Format(time.RFC3339), lt=endUTC.Format(time.RFC3339)
func DayRangeUTC(t time.Time, loc *time.Location) (startUTC time.Time, endUTC time.Time) {
	start, end := DayRange(t, loc)
	return start.UTC(), end.UTC()
}

// MonthRangeFrom returns [start, end) within the month that contains t, evaluated in loc.
//
// Rules:
//   - endExclusive is the start of the next month (00:00:00).
//   - if t is before "now" (current time), start = now (in loc).
//   - otherwise, start = start-of-day of t (in loc).
//
// Intended for range queries (gte/lt).
func MonthRangeFrom(t time.Time, loc *time.Location) (start time.Time, endExclusive time.Time) {
	if loc == nil {
		loc = MustVNLocation()
	}

	now := ToLocation(time.Now(), loc)
	tv := ToLocation(t, loc)

	// Determine start
	if tv.Before(now) {
		start = TimeOfDay(now, loc)
	} else {
		start = StartOfDay(tv, loc)
	}

	// End = start of next month
	year, month, _ := tv.Date()
	endExclusive = time.Date(year, month+1, 1, 0, 0, 0, 0, loc)
	endExclusive = endExclusive.Add(-time.Nanosecond)

	return
}

// IsBeforeTodayVN reports whether t is before today
// when both are evaluated in Asia/Ho_Chi_Minh timezone.
//
// Comparison is done at day granularity (start-of-day).
func IsBeforeTodayVN(t time.Time) bool {
	loc := MustVNLocation()

	tDay := StartOfDay(t, loc)
	today := StartOfDay(time.Now(), loc)

	return tDay.Before(today)
}
