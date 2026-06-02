package validate

import (
	"fmt"
	"strings"
	"time"
)

func IsTimeValid(t *time.Time) bool {
	return t != nil && !t.IsZero()
}

// InOnlineTimeRange checks whether now is within the online opening/closing range.
//
// Rules:
//   - both empty => always available
//   - only one empty => invalid config
//   - start == end => open 24 hours
//   - supports overnight range, e.g. 22:00 -> 06:00
//   - range is [start, end)
func InOnlineTimeRange(now time.Time, opening, closing string) (bool, error) {
	opening = strings.TrimSpace(opening)
	closing = strings.TrimSpace(closing)

	switch {
	case opening == "" && closing == "":
		return true, nil
	case opening == "" || closing == "":
		return false, fmt.Errorf("invalid online time range: opening=%q closing=%q", opening, closing)
	}

	startMinutes, err := parseHHMM(opening)
	if err != nil {
		return false, fmt.Errorf("invalid opening time: %w", err)
	}

	endMinutes, err := parseHHMM(closing)
	if err != nil {
		return false, fmt.Errorf("invalid closing time: %w", err)
	}

	currentMinutes := now.Hour()*60 + now.Minute()

	// Treat equal start/end as 24-hour open.
	if startMinutes == endMinutes {
		return true, nil
	}

	// Normal same-day range: 07:00 -> 22:00
	if startMinutes < endMinutes {
		return currentMinutes >= startMinutes && currentMinutes < endMinutes, nil
	}

	// Overnight range: 22:00 -> 06:00
	return currentMinutes >= startMinutes || currentMinutes < endMinutes, nil
}

func InOnlineTimeRangeBool(now time.Time, opening, closing string) bool {
	opening = strings.TrimSpace(opening)
	closing = strings.TrimSpace(closing)

	if opening == "" && closing == "" {
		return true
	}
	if opening == "" || closing == "" {
		return false
	}

	start, err := parseHHMM(opening)
	if err != nil {
		return false
	}

	end, err := parseHHMM(closing)
	if err != nil {
		return false
	}

	cur := now.Hour()*60 + now.Minute()

	if start == end {
		return true
	}
	if start < end {
		return cur >= start && cur < end
	}
	return cur >= start || cur < end
}

func parseHHMM(s string) (int, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, err
	}
	return t.Hour()*60 + t.Minute(), nil
}
