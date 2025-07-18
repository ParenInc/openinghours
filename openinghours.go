package openinghours

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// AmenityOpeningHours contains an amenity's opening and closing times within a given week. The
// values are in minutes since the beginning of the day.
//
// For example, an amenity that has OpeningHours of
//
//	[]OpeningHours {
//	    {
//	        Open: &TimeInWeek{ Weekday: time.Tuesday, MinutesSinceMidnight: 360 },
//	        Close: &TimeInWeek{ Weekday: time.Tuesday, MinutesSinceMidnight: 1200 },
//	    }, {
//	        Open: &TimeInWeek{ Weekday: time.Friday, MinutesSinceMidnight: 630 },
//	        Close: &TimeInWeek{ Weekday: time.Friday, MinutesSinceMidnight: 780 },
//	    },
//	}
//
// would mean that the amenity is open
// * tuesdays, from 06:00 (6am) to 20:00 (8pm); and
// * fridays, from 10:30 (10:30am) to 13:00 (1pm).
type OpeningHours struct {
	Open  *TimeInWeek
	Close *TimeInWeek
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
func GetHumanReadableTimes(s string) (map[string][]TimeRange, error) {
	strs := strings.Split(s, ",")

	openingTimes := make(map[string][]TimeRange)
	for _, str := range strs {
		if str == "" {
			continue
		}

		parts := strings.Split(str, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid opening hours string `%s`", str)
		}

		openingWeekInt, openingWeekday, openingTime, err := getHumanReadableTime(parts[0], false)
		if err != nil {
			return nil, fmt.Errorf("invalid opening hours string `%s`", str)
		}
		closingWeekInt, closingWeekday, closingTime, err := getHumanReadableTime(parts[1], true)
		if err != nil {
			return nil, fmt.Errorf("invalid opening hours string `%s`", str)
		}

		if openingWeekInt == closingWeekInt {
			addTimeToWeek(openingTimes, openingWeekday, openingTime, closingTime)
		} else {
			addTimeToWeek(openingTimes, openingWeekday, openingTime, "24:00")
			currentWeekInt := openingWeekInt
			currentWeekInt.Next()
			for currentWeekInt != closingWeekInt {
				addTimeToWeek(openingTimes, getWeekDay(currentWeekInt), "00:00", "24:00")
				currentWeekInt.Next()
			}
			addTimeToWeek(openingTimes, closingWeekday, "00:00", closingTime)
		}
	}
	return openingTimes, nil
}

func addTimeToWeek(times map[string][]TimeRange, weekday string, openingTime string, closingTime string) {
	times[weekday] = append(times[weekday], TimeRange{Open: openingTime, Close: closingTime})
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

	minutesSinceMidnight, err := parseMinutesSinceMidnight(matches[2], matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid time in `%s`: %s", v, err)
	}

	timeInWeek := TimeInWeek{
		Weekday:              weekday,
		MinutesSinceMidnight: minutesSinceMidnight,
	}

	return &timeInWeek, nil
}

func parseWeekDay(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil || (i < 1 || i > 7) {
		return 0
	}

	return i
}

func getWeekDay(weekday weekInt) string {
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

// GetWeekInt converts a weekday name (case-insensitive) to its corresponding RFC 3339 complient integer (1-7)
// or returns 0 for invalid input.
func GetWeekInt(v string) int {
	switch strings.ToLower(v) {
	case "monday":
		return 1
	case "tuesday":
		return 2
	case "wednesday":
		return 3
	case "thursday":
		return 4
	case "friday":
		return 5
	case "saturday":
		return 6
	case "sunday":
		return 7
	default:
		return 0
	}
}

// getHumanReadableTime parses a time string and returns the corresponding weekday, time, and week integer.
// The `v` parameter is the time string to parse, which should follow the pattern "W{week}T{hour}:{minute}:{second}".
// The `endTime` flag determines if the time at the end of a period should adjust to "24:00" of previous weekday if "00:00".
// Returns the week integer, friendly weekday name, formatted time string, and an error if parsing fails.
func getHumanReadableTime(v string, endTime bool) (weekInt, string, string, error) {
	re := regexp.MustCompile(`^W(\d)T(\d{2}):(\d{2}):\d{2}$`)
	matches := re.FindStringSubmatch(v)
	if len(matches) < 2 {
		return 0, "", "", fmt.Errorf("invalid value `%s`", v)
	}
	weekIndex, _ := strconv.Atoi(matches[1])
	weekI := weekInt(weekIndex)
	hours := matches[2]
	minutes := matches[3]
	if hours == "00" && minutes == "00" && endTime {
		weekI.Previous()
		hours = "24"
	}
	weekday := getWeekDay(weekI)
	time := fmt.Sprintf("%s:%s", hours, minutes)

	return weekI, weekday, time, nil
}

func parseMinutesSinceMidnight(v1, v2 string) (int, error) {
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

type weekInt int

func (w *weekInt) Next() {
	*w = *w%7 + 1
}

func (w *weekInt) Previous() {
	*w = *w - 1
	if *w < 1 {
		*w = 7
	}
}
