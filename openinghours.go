package openinghours

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// OpeningHours contains an opening and closing times within a given week. The
// values are in minutes since the beginning of the day.
//
// For example, a location that has OpeningHours of
//
//	[]OpeningHours {
//	    {
//	        Open: &timeInWeek{ weekday: time.Tuesday, minutesSinceMidnight: 360 },
//	        Close: &timeInWeek{ weekday: time.Tuesday, minutesSinceMidnight: 1200 },
//	    }, {
//	        Open: &timeInWeek{ weekday: time.Friday, minutesSinceMidnight: 630 },
//	        Close: &timeInWeek{ weekday: time.Friday, minutesSinceMidnight: 780 },
//	    },
//	}
//
// would mean that it is open
// * tuesdays, from 06:00 (6am) to 20:00 (8pm); and
// * fridays, from 10:30 (10:30am) to 13:00 (1pm).

const (
	TwentyFourSevenString = "W1T00:00:00/W7T24:00:00"
)

var (
	TwentyFourSevenOH = OpeningHours{
		Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 0},
		Close: &TimeInWeek{Weekday: 7, MinutesSinceMidnight: 1440},
	}
)

type OpeningHours struct {
	Open  *TimeInWeek
	Close *TimeInWeek
}

// String returns the opening hours of the amenity in a string. Unfortunately, there are no formats
// in either standards (RFC 3339 or ISO 8601) to represent a recurring time within a given week, so
// one is invented here.
//
// Using the same example as above, the resulting strings would be
// * "W2T06:00:00/W2T20:00:00"; and
// * "W5T10:30:00/W5T13:00:00".
//
// Contrary to the stdlib's time, the start of the week is monday, to follow RFC 3339.
//
// No time zone information is provided, as the opening hours are static within the given day, ie.
// they don't change during a daylight saving time change.
func (oh OpeningHours) String() string {
	var open string
	if oh.Open != nil {
		open = fmt.Sprintf(
			"W%dT%02d:%02d:00",
			oh.Open.Weekday,
			oh.Open.MinutesSinceMidnight/60,
			oh.Open.MinutesSinceMidnight%60,
		)
	}

	var close string
	if oh.Close != nil {
		close = fmt.Sprintf(
			"W%dT%02d:%02d:00",
			oh.Close.Weekday,
			oh.Close.MinutesSinceMidnight/60,
			oh.Close.MinutesSinceMidnight%60,
		)
	}

	return fmt.Sprintf("%s/%s", open, close)
}

// OpeningHoursSliceToString converts a slice of OpeningHours into a single string representation like "W1T08:00:00/W1T16:00:00,W2T06:00:00/W2T20:00:00".
func OpeningHoursSliceToString(ohs []OpeningHours) string {
	openingHoursStr := make([]string, len(ohs))
	for i, openingHours := range ohs {
		openingHoursStr[i] = openingHours.String()
	}

	return strings.Join(openingHoursStr, ",")
}

// ParseOpeningHours does the opposite of OpeningHours.String method. It converts a string like
// "W0T08:00:00/W0T20:00:00" into a []OpeningHours.
func ParseOpeningHours(v string) ([]OpeningHours, error) {
	strs := strings.Split(v, ",")

	ohs := make([]OpeningHours, 0, len(strs))
	for _, str := range strs {
		if str == "" {
			continue
		}

		parts := strings.Split(str, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid opening hours string `%s`", str)
		}

		openingHours, err := parseTimeInWeek(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid opening hours: %s", err)
		}

		closingHours, err := parseTimeInWeek(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid closing hours: %s", err)
		}

		oh := OpeningHours{
			Open:  openingHours,
			Close: closingHours,
		}

		ohs = append(ohs, oh)
	}

	return ohs, nil
}

type TimeRange struct {
	Open  string `json:"open"`
	Close string `json:"close"`
}

// GetHumanReadableTimes Returns a map of weekdays and timeRanges for the given opening hours string.
// The opening hours string is a comma-separated list of opening and closing times, each of which is
// formatted as "W{week}T{hour}:{minute}:{second}".
// For example, given the opening hours string "W3T10:00:00/W3T20:30:00,W5T10:00:00/W5T12:00:00,W5T13:00:00/W5T21:00:00",
// the returned map would be:
//
//	openingHours: {
//	 Monday: nil
//	 Tuesday: nil
//	 Wednesday: [{open: "10:00", close: "20:30"}]
//	 Thursday: nil
//	 Friday: [{open: "10:00", close: "12:00"}, {open: "13:00", close: "21:00"}]
//	 ...
//	}
func GetHumanReadableTimes(ohs []OpeningHours) map[string][]TimeRange {
	if len(ohs) == 0 {
		return nil
	}

	openingTimes := make(map[string][]TimeRange)
	for _, oh := range ohs {
		if oh.Close.MinutesSinceMidnight == 0 {
			setPreviousDay(&oh.Close.Weekday)
			oh.Close.MinutesSinceMidnight = 1440 // 24:00
		}
		if oh.Open.Weekday == oh.Close.Weekday {
			addTimeToWeek(openingTimes, getWeekDay(oh.Open.Weekday), minutesSinceMidnightToTime(oh.Open.MinutesSinceMidnight), minutesSinceMidnightToTime(oh.Close.MinutesSinceMidnight))
		} else {
			addTimeToWeek(openingTimes, getWeekDay(oh.Open.Weekday), minutesSinceMidnightToTime(oh.Open.MinutesSinceMidnight), "24:00")
			setNextDay(&oh.Open.Weekday)
			for oh.Open.Weekday != oh.Close.Weekday {
				addTimeToWeek(openingTimes, getWeekDay(oh.Open.Weekday), "00:00", "24:00")
				setNextDay(&oh.Open.Weekday)
			}
			addTimeToWeek(openingTimes, getWeekDay(oh.Close.Weekday), "00:00", minutesSinceMidnightToTime(oh.Close.MinutesSinceMidnight))
		}
	}
	return openingTimes
}

func addTimeToWeek(times map[string][]TimeRange, weekday string, openingTime string, closingTime string) {
	times[weekday] = append(times[weekday], TimeRange{Open: openingTime, Close: closingTime})
}

type OCPIOpeningTimes struct {
	TwentyFourSeven bool                `json:"twentyfourseven" example:"false"`
	RegularHours    *[]OCPIRegularHours `json:"regular_hours,omitempty"`
}

type OCPIRegularHours struct {
	Weekday     int    `json:"weekday" example:"1"`
	PeriodBegin string `json:"period_begin" example:"06:00"`
	PeriodEnd   string `json:"period_end" example:"22:00"` //  Must be later than period_begin or be "00:00" to signal that the charging station is open until midnight at the end of the day.
}

// GetOCPIOpeningTimes converts a slice of OpeningHours into an OCPIOpeningTimes struct.
// If the opening hours are 24/7, it returns an OCPIOpeningTimes with TwentyFourSeven set to true.
// Example:
//   ohs := []OpeningHours{
//       {Open: &TimeInWeek{Weekday: 1, minutesSinceMidnight: 0}, Close: &TimeInWeek{Weekday: 7, minutesSinceMidnight: 1440}},
//   }
//   ocpiOpeningTimes := GetOCPIOpeningTimes(ohs)
//   // ocpiOpeningTimes will be OCPIOpeningTimes{TwentyFourSeven: true}
//
// If the opening hours are not 24/7, it returns an OCPIOpeningTimes with
// RegularHours containing the opening and closing times for each day of the week.
// Example:
//   ohs := []OpeningHours{
//       {Open: &TimeInWeek{Weekday: 1, minutesSinceMidnight: 360}, Close: &TimeInWeek{Weekday: 1, minutesSinceMidnight: 1200}},
//       {Open: &TimeInWeek{Weekday: 5, minutesSinceMidnight: 630}, Close: &TimeInWeek{Weekday: 5, minutesSinceMidnight: 780}},
//   }
//   ocpiOpeningTimes := GetOCPIOpeningTimes(ohs)
//   // ocpiOpeningTimes will be OCPIOpeningTimes{
//       TwentyFourSeven: false,
//       RegularHours: &[]OCPIRegularHours{
//           {Weekday: 1, PeriodBegin: "06:00", PeriodEnd: "20:00"},
//           {Weekday: 5, PeriodBegin: "10:30", PeriodEnd: "13:00"},
//       },
//   }
//

func GetOCPIOpeningTimes(ohs []OpeningHours) OCPIOpeningTimes {
	if isTwentyFourSeven(ohs) {
		return OCPIOpeningTimes{TwentyFourSeven: true}
	}

	var regularHours []OCPIRegularHours
	for _, oh := range ohs {
		switch oh.Close.MinutesSinceMidnight {
		case 0:
			setPreviousDay(&oh.Close.Weekday)
		case 1440:
			oh.Close.MinutesSinceMidnight = 0 // 24:00 is represented as 00:00 in the OCPI spec
		}

		if oh.Open.Weekday == oh.Close.Weekday {
			regularHours = append(regularHours, OCPIRegularHours{
				Weekday:     oh.Open.Weekday,
				PeriodBegin: minutesSinceMidnightToTime(oh.Open.MinutesSinceMidnight),
				PeriodEnd:   minutesSinceMidnightToTime(oh.Close.MinutesSinceMidnight),
			})
			continue
		} else {
			regularHours = append(regularHours, OCPIRegularHours{
				Weekday:     oh.Open.Weekday,
				PeriodBegin: minutesSinceMidnightToTime(oh.Open.MinutesSinceMidnight),
				PeriodEnd:   "00:00",
			})
			setNextDay(&oh.Open.Weekday)
			for oh.Open.Weekday != oh.Close.Weekday {
				regularHours = append(regularHours, OCPIRegularHours{
					Weekday:     oh.Open.Weekday,
					PeriodBegin: "00:00",
					PeriodEnd:   "00:00",
				})
				setNextDay(&oh.Open.Weekday)
			}
			regularHours = append(regularHours, OCPIRegularHours{
				Weekday:     oh.Close.Weekday,
				PeriodBegin: "00:00",
				PeriodEnd:   minutesSinceMidnightToTime(oh.Close.MinutesSinceMidnight),
			})
		}
	}
	if len(regularHours) == 0 {
		return OCPIOpeningTimes{}
	}

	return OCPIOpeningTimes{
		TwentyFourSeven: false,
		RegularHours:    &regularHours,
	}
}

// ParseStringWeekdayToTimeWeekday converts a string representation of a weekday
// (e.g., "monday", "tuesday") to the corresponding int value.
func ParseStringWeekdayToTimeWeekday(dayStr string) (int, error) {
	switch strings.ToLower(dayStr) {
	case "monday", "mon":
		return 1, nil
	case "tuesday", "tue":
		return 2, nil
	case "wednesday", "wed":
		return 3, nil
	case "thursday", "thu":
		return 4, nil
	case "friday", "fri":
		return 5, nil
	case "saturday", "sat":
		return 6, nil
	case "sunday", "sun":
		return 7, nil
	default:
		return 0, fmt.Errorf("invalid weekday: %s", dayStr)
	}
}

// ParseMinutesSinceMidnight parses hours and minutes strings into total minutes since midnight.
// e.g. ("08", "30") -> 510
func ParseMinutesSinceMidnight(v1, v2 string) (int, error) {
	hours, err := strconv.Atoi(v1)
	if err != nil || (hours < 0 || hours > 24) {
		return 0, fmt.Errorf("invalid hours value")
	}

	minutes, err := strconv.Atoi(v2)
	if err != nil || (minutes < 0 || minutes > 59) {
		return 0, fmt.Errorf("invalid minutes value")
	}

	if hours == 24 && minutes != 0 {
		return 0, fmt.Errorf("invalid value")
	}

	return hours*60 + minutes, nil
}

func isTwentyFourSeven(ohs []OpeningHours) bool {
	if len(ohs) == 0 {
		return false
	}

	for _, oh := range ohs {
		if oh.Open == nil || oh.Close == nil {
			return false
		}
		if oh.Open.Weekday != 1 || oh.Close.Weekday != 7 {
			return false
		}
		if oh.Open.MinutesSinceMidnight != 0 || oh.Close.MinutesSinceMidnight != 1440 {
			return false
		}
	}

	return true
}

// TimeInWeek contains a time within the week, given by the weekday number and the minutes since
// midnight.
//
// Note that the Weekday is as per RFC 3339, not stdlib's time.Weekday.
type TimeInWeek struct {
	Weekday              int
	MinutesSinceMidnight int
}

func parseTimeInWeek(v string) (*TimeInWeek, error) {
	if v == "" {
		return nil, nil
	}

	re := regexp.MustCompile(`^W(\d)T(\d{2}):(\d{2}):\d{2}$`)
	matches := re.FindStringSubmatch(v)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid value `%s`", v)
	}

	weekday := parseWeekDay(matches[1])
	if weekday == 0 {
		return nil, fmt.Errorf("invalid workday in `%s`: expected to be between 1 (monday) and 7 (sunday)", v)
	}

	minutesSinceMidnight, err := ParseMinutesSinceMidnight(matches[2], matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid time in `%s`: %s", v, err)
	}

	tiw := TimeInWeek{
		Weekday:              weekday,
		MinutesSinceMidnight: minutesSinceMidnight,
	}

	return &tiw, nil
}

func parseWeekDay(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil || (i < 1 || i > 7) {
		return 0
	}

	return i
}

func getWeekDay(weekday int) string {
	switch weekday {
	case 1:
		return "monday"
	case 2:
		return "tuesday"
	case 3:
		return "wednesday"
	case 4:
		return "thursday"
	case 5:
		return "friday"
	case 6:
		return "saturday"
	case 7:
		return "sunday"
	default:
		return ""
	}
}

func minutesSinceMidnightToTime(minutesSinceMidnight int) string {
	hours := minutesSinceMidnight / 60
	minutes := minutesSinceMidnight % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func setNextDay(w *int) {
	*w = *w%7 + 1
}

func setPreviousDay(w *int) {
	*w = *w - 1
	if *w < 1 {
		*w = 7
	}
}
