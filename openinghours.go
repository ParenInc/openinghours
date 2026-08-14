package openinghours

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

// freeformParser attempts to parse a free-form, human-written opening hours string. It returns
// the parsed OpeningHours and true if the string matched its format, or false if the string
// doesn't match, so that ParseFreeformOpeningHours can try the next parser in the list.
type freeformParser func(string) ([]OpeningHours, bool, error)

// freeformParsers is the list of supported free-form string formats, tried in order by
// ParseFreeformOpeningHours. Add new formats by appending a parser here.
var freeformParsers = []freeformParser{
	parseTwentyFourSevenDaily,
	parseMultiDaySegments,
	parseDailyRange,
}

// parentheticalRegexp matches a parenthetical aside, eg. "(May-Sept)" or "(Sunday hours in winter
// only)". These are treated as unstructured notes and stripped before parsing, the same way
// free-text clauses like "for library visitors" are ignored.
var parentheticalRegexp = regexp.MustCompile(`\([^)]*\)`)

// ParseFreeformOpeningHours converts a free-form, human-written opening hours string (e.g.
// "9am-9pm") into a []OpeningHours, trying each parser in freeformParsers in turn and returning
// the result of the first one that matches. Parenthetical asides (eg. "(May-Sept)") are stripped
// before parsing.
func ParseFreeformOpeningHours(input string) ([]OpeningHours, error) {
	trimmed := strings.TrimSpace(parentheticalRegexp.ReplaceAllString(input, ""))
	if trimmed == "" {
		return nil, fmt.Errorf("empty opening hours string")
	}

	for _, parse := range freeformParsers {
		ohs, ok, err := parse(trimmed)
		if err != nil {
			return nil, err
		}
		if ok {
			return ohs, nil
		}
	}

	return nil, fmt.Errorf("unrecognized opening hours format: %q", input)
}

// FreeformOpeningHoursToISOString converts a free-form opening hours string (e.g. "9am-9pm")
// directly into its ISO-like string representation, as produced by OpeningHoursSliceToString.
func FreeformOpeningHoursToISOString(input string) (string, error) {
	ohs, err := ParseFreeformOpeningHours(input)
	if err != nil {
		return "", err
	}

	return OpeningHoursSliceToString(ohs), nil
}

var twentyFourSevenDailyRegexp = regexp.MustCompile(`(?i)^24\s*(?:hours?\s*daily|/\s*7)$`)

// parseTwentyFourSevenDaily matches a string like "24 hours daily", "24 Hours Daily" or "24/7",
// meaning the amenity is open every hour of every day of the week.
func parseTwentyFourSevenDaily(input string) ([]OpeningHours, bool, error) {
	if !twentyFourSevenDailyRegexp.MatchString(input) {
		return nil, false, nil
	}

	return []OpeningHours{TwentyFourSevenOH}, true, nil
}

// weekdayAbbreviations maps the various abbreviated forms of a weekday name, as commonly found in
// human-written opening hours (eg. "M", "Th", "Mon", "Sat"), to their RFC 3339 weekday number. A
// bare "T" always means Tuesday, since Thursday is conventionally written out as "Th" in this
// abbreviation style, leaving "T" unambiguous in practice.
var weekdayAbbreviations = map[string]int{
	"m":         1,
	"mon":       1,
	"monday":    1,
	"t":         2,
	"tu":        2,
	"tue":       2,
	"tues":      2,
	"tuesday":   2,
	"w":         3,
	"wed":       3,
	"wednesday": 3,
	"th":        4,
	"thu":       4,
	"thur":      4,
	"thurs":     4,
	"thursday":  4,
	"f":         5,
	"fri":       5,
	"friday":    5,
	"sa":        6,
	"sat":       6,
	"saturday":  6,
	"su":        7,
	"sun":       7,
	"sunday":    7,
}

// parseWeekdayAbbreviation converts an abbreviated weekday string (see weekdayAbbreviations) into
// its RFC 3339 weekday number.
func parseWeekdayAbbreviation(v string) (int, error) {
	weekday, ok := weekdayAbbreviations[strings.ToLower(v)]
	if !ok {
		return 0, fmt.Errorf("invalid weekday abbreviation: %q", v)
	}

	return weekday, nil
}

// expandWeekdayRange returns every weekday number from `from` to `to`, inclusive, wrapping around
// the week (eg. 6 (saturday) to 1 (monday) => [6, 7, 1]) if `to` comes before `from`.
func expandWeekdayRange(from, to int) []int {
	weekdays := []int{from}
	for weekdays[len(weekdays)-1] != to {
		next := weekdays[len(weekdays)-1]
		setNextDay(&next)
		weekdays = append(weekdays, next)
	}

	return weekdays
}

// daySegmentRegexp matches a single "<time range> <day spec>" clause, eg. "9am-9pm M-Th",
// "9am-6pm F-Sat", "1pm-5pm Sun", "10am-5pm daily", "12pm-5pm, Mon-Sat" (a comma is accepted in
// place of, or alongside, the whitespace between the time range and the day spec), "5am to 10pm
// daily" (the word "to" is accepted as an alternative to the hyphen), "9-5pm Sat" (the opening
// time's meridiem is optional and defaults to "am" when omitted - the closing time's meridiem is
// still required) or "11am-6pm W-F and Sun" / "10am-6pm M & W" (a day spec may itself be a
// ","/"and"/"&"-joined list of days or day ranges). The day spec is captured permissively
// (anything at all) because resolveDaySpec knows how to recover a valid day spec from a trailing
// free-text note glued onto it without a proper separator (eg. "M-F daily", "daily for guests").
var daySegmentRegexp = regexp.MustCompile(`(?i)^(\d{1,2})(?::([0-5][0-9]))?(?:\s*(am|pm))?(?:\s*-\s*|\s+to\s+)(\d{1,2})(?::([0-5][0-9]))?\s*(am|pm)[,\s]+(.+)$`)

// daySpecSeparatorRegexp splits a day spec into its individual day tokens or day ranges, eg.
// "W-F and Sun" splits into ["W-F", "Sun"], "Mon, Wed and Fri" splits into ["Mon", "Wed", "Fri"],
// and "M & W" splits into ["M", "W"].
var daySpecSeparatorRegexp = regexp.MustCompile(`(?i)\s*,\s*|\s+and\s+|\s*&\s*`)

// hasWordPrefix reports whether runes begins with word, case-insensitively, as a whole word (ie.
// not immediately followed by another letter, so "closed" matches but "closedown" wouldn't).
func hasWordPrefix(runes []rune, word string) bool {
	if len(runes) < len(word) || !strings.EqualFold(string(runes[:len(word)]), word) {
		return false
	}

	return len(runes) == len(word) || !unicode.IsLetter(runes[len(word)])
}

// splitOnComma splits a string on commas, but only when a comma is followed (after any
// whitespace) by a digit - the start of a new time range - or by the word "closed" - the start of
// a new exclusion clause, eg. "closed Sun". A comma directly between a time range and its own day
// spec (eg. "12pm-5pm, Mon-Sat") is left alone, so daySegmentRegexp can match it as one clause.
// RE2 (Go's regexp package) has no lookahead, so this can't be expressed as a single regexp.
func splitOnComma(input string) []string {
	runes := []rune(input)

	var segments []string
	start := 0
	for i, r := range runes {
		if r != ',' {
			continue
		}

		j := i + 1
		for j < len(runes) && unicode.IsSpace(runes[j]) {
			j++
		}
		if j < len(runes) && (unicode.IsDigit(runes[j]) || hasWordPrefix(runes[j:], "closed")) {
			segments = append(segments, string(runes[start:i]))
			start = i + 1
		}
	}
	segments = append(segments, string(runes[start:]))

	return segments
}

// allWeekdays are the RFC 3339 weekday numbers for every day of the week, monday through sunday.
var allWeekdays = []int{1, 2, 3, 4, 5, 6, 7}

// weekdaysMonToFri are the RFC 3339 weekday numbers for the "weekdays" keyword.
var weekdaysMonToFri = []int{1, 2, 3, 4, 5}

// resolveDaySpec resolves a day spec (the day portion of a daySegmentRegexp match) into the
// weekday numbers it covers, or false if nothing in it resolves to a recognizable day at all. A
// day spec is the keyword "daily"/"every day" (every day of the week), "weekdays"/"weekday"
// (monday through friday), or one or more day tokens/day ranges (eg. "Mon", "W-F") joined by
// ","/"and"/"&" (eg. "W-F and Sun", "Mon, Wed and Fri", "M & W").
//
// Each token is resolved independently via resolveDaySpecToken, which progressively drops
// trailing words and retries if a token doesn't resolve as-is. This recovers the recoverable days
// from a free-text note glued onto (or alongside) a valid day spec with no clean separator, eg.:
//   - "M-F daily" and "daily for guests" (a single token, with no list separator at all) recover
//     "M-F" and "daily" respectively;
//   - "Sat, for client use only" and "M-F, for employee use only" (a real day followed by a
//     comma-glued note) recover "Sat" and "M-F", silently dropping the unresolvable second token;
//   - "daily for staff, students and faculty with a permit" recovers "daily" from its first
//     token, dropping "students" and "faculty with a permit".
//
// A token that never resolves to anything, even after trimming down to a single word (eg. "Xyz",
// "students"), is silently dropped rather than treated as an error: free-text notes are common and
// varied enough that erroring on an unrecognized word is more likely to reject a legitimate note
// than to catch a genuine typo in a day list. If nothing in the whole day spec resolves to any
// weekday at all, ok is false, so the caller can treat the whole clause as unmatched, the same as
// if it never looked like a day spec in the first place.
//
// Weekdays are returned in first-seen order, deduplicated, so an accidental overlap between
// tokens doesn't produce duplicate OpeningHours entries.
func resolveDaySpec(daySpec string) ([]int, bool) {
	var weekdays []int
	seen := make(map[int]bool, 7)
	for _, token := range daySpecSeparatorRegexp.Split(daySpec, -1) {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		tokenWeekdays, ok := resolveDaySpecToken(token)
		if !ok {
			continue
		}

		for _, weekday := range tokenWeekdays {
			if seen[weekday] {
				continue
			}

			seen[weekday] = true
			weekdays = append(weekdays, weekday)
		}
	}

	return weekdays, len(weekdays) > 0
}

// resolveDaySpecToken resolves one token - a single day, day range, or recognized keyword - into
// the weekday numbers it covers, progressively dropping trailing words and retrying if it doesn't
// resolve as-is (see resolveDaySpec).
func resolveDaySpecToken(token string) ([]int, bool) {
	if weekdays, ok := resolveSingleDaySpec(token); ok {
		return weekdays, true
	}

	words := strings.Fields(token)
	for n := len(words) - 1; n >= 1; n-- {
		if weekdays, ok := resolveSingleDaySpec(strings.Join(words[:n], " ")); ok {
			return weekdays, true
		}
	}

	return nil, false
}

// resolveSingleDaySpec resolves one day token, day range, or recognized keyword ("daily", "every
// day", "weekdays", "weekday") into the weekday numbers it covers. ok is false if s doesn't
// resolve to anything recognizable.
func resolveSingleDaySpec(s string) ([]int, bool) {
	switch {
	case strings.EqualFold(s, "daily"), strings.EqualFold(s, "every day"):
		return allWeekdays, true
	case strings.EqualFold(s, "weekdays"), strings.EqualFold(s, "weekday"):
		return weekdaysMonToFri, true
	}

	fromDay, toDay, err := parseDayOrDayRange(s)
	if err != nil {
		return nil, false
	}

	return expandWeekdayRange(fromDay, toDay), true
}

// parseDayOrDayRange parses a single day token (eg. "Sat") or day range (eg. "M-Th") into its
// start and end weekday numbers. A single day's start and end are the same weekday. "S-S" is
// special-cased to mean "Saturday-Sunday", a common shorthand; a bare "S" outside of that idiom is
// genuinely ambiguous between Saturday and Sunday and isn't otherwise supported.
func parseDayOrDayRange(token string) (int, int, error) {
	if strings.EqualFold(token, "s-s") {
		return 6, 7, nil
	}

	parts := strings.SplitN(token, "-", 2)

	fromDay, err := parseWeekdayAbbreviation(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}

	if len(parts) == 1 {
		return fromDay, fromDay, nil
	}

	toDay, err := parseWeekdayAbbreviation(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}

	return fromDay, toDay, nil
}

// parseTimeRangeMinutes parses the "<time>-<time>" (or "<time> to <time>") prefix shared by
// dailyRangeRegexp and daySegmentRegexp (whose first six capture groups have the same shape) into
// opening and closing minutes since midnight. overnight is true when the closing time is earlier
// in the day than the opening time, meaning the range spans past midnight into the next day (eg.
// "10pm-6am" opens at 22:00 and closes at 06:00 the following day).
func parseTimeRangeMinutes(matches []string) (openMinutes, closeMinutes int, overnight bool, err error) {
	openMinutes, err = parse12HourTime(matches[1], matches[2], matches[3])
	if err != nil {
		return 0, 0, false, fmt.Errorf("invalid opening time in %q: %s", matches[0], err)
	}

	closeMinutes, err = parse12HourClosingTime(matches[4], matches[5], matches[6])
	if err != nil {
		return 0, 0, false, fmt.Errorf("invalid closing time in %q: %s", matches[0], err)
	}

	if closeMinutes == openMinutes {
		return 0, 0, false, fmt.Errorf("closing time must be after opening time in %q", matches[0])
	}

	return openMinutes, closeMinutes, closeMinutes < openMinutes, nil
}

// buildOpeningHour builds a single OpeningHours entry for the given weekday and time window,
// rolling the closing weekday forward to the next day when overnight is true.
func buildOpeningHour(weekday, openMinutes, closeMinutes int, overnight bool) OpeningHours {
	closeWeekday := weekday
	if overnight {
		setNextDay(&closeWeekday)
	}

	return OpeningHours{
		Open:  &TimeInWeek{Weekday: weekday, MinutesSinceMidnight: openMinutes},
		Close: &TimeInWeek{Weekday: closeWeekday, MinutesSinceMidnight: closeMinutes},
	}
}

// buildOpeningHoursForWeekdays builds one OpeningHours entry per weekday for a single time
// window.
func buildOpeningHoursForWeekdays(weekdays []int, openMinutes, closeMinutes int, overnight bool) []OpeningHours {
	ohs := make([]OpeningHours, 0, len(weekdays))
	for _, weekday := range weekdays {
		ohs = append(ohs, buildOpeningHour(weekday, openMinutes, closeMinutes, overnight))
	}

	return ohs
}

// buildEveryDayOpeningHours builds one OpeningHours entry per weekday from a bare "<time>-<time>"
// match (from dailyRangeRegexp).
func buildEveryDayOpeningHours(matches []string) ([]OpeningHours, error) {
	openMinutes, closeMinutes, overnight, err := parseTimeRangeMinutes(matches)
	if err != nil {
		return nil, err
	}

	return buildOpeningHoursForWeekdays(allWeekdays, openMinutes, closeMinutes, overnight), nil
}

// allBareTimeRanges reports whether every segment is a bare time range with no day spec of its
// own (eg. "12am-7am"), as matched by dailyRangeRegexp. An empty list reports false, since there's
// nothing to share a trailing day spec with.
func allBareTimeRanges(segments []string) bool {
	if len(segments) == 0 {
		return false
	}

	for _, segment := range segments {
		if !dailyRangeRegexp.MatchString(segment) {
			return false
		}
	}

	return true
}

// timeWindow is an opening/closing time pair, in minutes since midnight, plus whether it spans
// past midnight into the next day.
type timeWindow struct {
	openMinutes, closeMinutes int
	overnight                 bool
}

// parseSharedDaySpecGroup parses a list of segments whose final segment carries the only day
// spec, which applies to every time range in the list, eg. ["12am-7am", "9am-5:30pm",
// "7:30pm-12am daily"] (all three windows apply daily). Every segment but the last must match
// dailyRangeRegexp, and the last must match daySegmentRegexp.
func parseSharedDaySpecGroup(segments []string) ([]OpeningHours, error) {
	lastMatches := daySegmentRegexp.FindStringSubmatch(segments[len(segments)-1])

	weekdays, ok := resolveDaySpec(lastMatches[7])
	if !ok {
		return nil, nil
	}

	windows := make([]timeWindow, 0, len(segments))
	for _, segment := range segments[:len(segments)-1] {
		openMinutes, closeMinutes, overnight, err := parseTimeRangeMinutes(dailyRangeRegexp.FindStringSubmatch(segment))
		if err != nil {
			return nil, err
		}

		windows = append(windows, timeWindow{openMinutes, closeMinutes, overnight})
	}

	lastOpenMinutes, lastCloseMinutes, lastOvernight, err := parseTimeRangeMinutes(lastMatches)
	if err != nil {
		return nil, err
	}
	windows = append(windows, timeWindow{lastOpenMinutes, lastCloseMinutes, lastOvernight})

	var ohs []OpeningHours
	for _, weekday := range weekdays {
		for _, w := range windows {
			ohs = append(ohs, buildOpeningHour(weekday, w.openMinutes, w.closeMinutes, w.overnight))
		}
	}

	return ohs, nil
}

// windowSeparatorRegexp splits a segment on " and " or "&" when it's being used to join multiple
// bare time ranges that share one trailing day spec, eg. "8am-12pm and 1pm-5pm M-F" or
// "6:30am-6pm & 8pm-10pm daily".
var windowSeparatorRegexp = regexp.MustCompile(`(?i)\s+and\s+|\s*&\s*`)

// parseAndSharedWindowSegment matches a segment made of two or more "and"/"&"-joined time ranges
// that share a single trailing day spec, eg. "8am-12pm and 1pm-5pm M-F" or "6:30am-6pm & 8pm-10pm
// daily". It reports matched=false (without error) if the segment doesn't look like this pattern -
// eg. a single time range whose day spec happens to contain "and"/"&" itself, like "11am-6pm W-F
// and Sun" or "10am-6pm M & W" - so the caller can fall back to treating it as an ordinary single
// "<time range> <day spec>" clause.
func parseAndSharedWindowSegment(segment string) ([]OpeningHours, bool, error) {
	parts := windowSeparatorRegexp.Split(segment, -1)
	if len(parts) < 2 || !daySegmentRegexp.MatchString(parts[len(parts)-1]) || !allBareTimeRanges(parts[:len(parts)-1]) {
		return nil, false, nil
	}

	ohs, err := parseSharedDaySpecGroup(parts)
	if err != nil {
		return nil, true, err
	}

	return ohs, true, nil
}

// twentyFourHoursDaySpecRegexp matches "24 hours <day spec>" or "24 hour <day spec>", eg. "24
// hours Sat-Sun", meaning the amenity is open the full day (00:00-24:00) on the given days. This
// is distinct from parseTwentyFourSevenDaily, which matches the keyword applied to the *whole*
// input string; this one is a single clause that can sit alongside other clauses in the same
// string, eg. "3pm-7:45am M-F, 24 hours Sat-Sun".
var twentyFourHoursDaySpecRegexp = regexp.MustCompile(`(?i)^24\s*hours?\s+(.+)$`)

// parseTwentyFourHoursDaySpecSegment matches a "24 hours <day spec>" segment.
func parseTwentyFourHoursDaySpecSegment(segment string) ([]OpeningHours, bool, error) {
	matches := twentyFourHoursDaySpecRegexp.FindStringSubmatch(segment)
	if matches == nil {
		return nil, false, nil
	}

	weekdays, ok := resolveDaySpec(matches[1])
	if !ok {
		return nil, false, nil
	}

	return buildOpeningHoursForWeekdays(weekdays, 0, 1440, false), true, nil
}

// parseSingleDaySegment matches an ordinary single "<time range> <day spec>" clause.
func parseSingleDaySegment(segment string) ([]OpeningHours, bool, error) {
	matches := daySegmentRegexp.FindStringSubmatch(segment)
	if matches == nil {
		return nil, false, nil
	}

	openMinutes, closeMinutes, overnight, err := parseTimeRangeMinutes(matches)
	if err != nil {
		return nil, true, err
	}

	weekdays, ok := resolveDaySpec(matches[7])
	if !ok {
		return nil, false, nil
	}

	return buildOpeningHoursForWeekdays(weekdays, openMinutes, closeMinutes, overnight), true, nil
}

// segmentParsers is the list of patterns tried, in order, for each individual comma-separated
// segment within a parseDaySegmentGroup group.
var segmentParsers = []func(string) ([]OpeningHours, bool, error){
	parseTwentyFourHoursDaySpecSegment,
	parseAndSharedWindowSegment,
	parseSingleDaySegment,
}

// parseDaySegmentGroup parses one semicolon-delimited group of one or more comma-separated
// segments. Normally each segment carries its own day spec, eg. "9am-9pm M-Th, 9am-6pm F". But:
//   - when a group's last segment is the only one with a day spec, and every segment before it is
//     a bare time range (eg. "12am-7am, 9am-5:30pm, 7:30pm-12am daily"), that day spec is shared
//     across all of them;
//   - a lone segment with no day spec at all, but a valid bare time range (eg. "8am-5pm" with
//     nothing else in the group), is treated as applying every day, same as parseDailyRange;
//   - "24 hours <day spec>" and "<time> and <time> ... <day spec>" segments are recognized by
//     segmentParsers.
//
// Segments that match none of these, such as free-text notes ("for library visitors", "2 hour
// charging limit") or explicit exclusions ("closed Sun"), are ignored.
func parseDaySegmentGroup(group string) ([]OpeningHours, error) {
	var trimmed []string
	for _, segment := range splitOnComma(group) {
		if segment = strings.TrimSpace(segment); segment != "" {
			trimmed = append(trimmed, segment)
		}
	}

	if len(trimmed) == 0 {
		return nil, nil
	}

	if len(trimmed) > 1 && daySegmentRegexp.MatchString(trimmed[len(trimmed)-1]) && allBareTimeRanges(trimmed[:len(trimmed)-1]) {
		return parseSharedDaySpecGroup(trimmed)
	}

	if len(trimmed) == 1 {
		if matches := dailyRangeRegexp.FindStringSubmatch(trimmed[0]); matches != nil {
			return buildEveryDayOpeningHours(matches)
		}
	}

	var ohs []OpeningHours
	for _, segment := range trimmed {
		for _, parse := range segmentParsers {
			segmentOhs, matched, err := parse(segment)
			if err != nil {
				return nil, err
			}
			if matched {
				ohs = append(ohs, segmentOhs...)
				break
			}
		}
	}

	return ohs, nil
}

// parseMultiDaySegments matches opening-hours strings made up of one or more semicolon-separated
// groups, eg. "9am-9pm M-Th, 9am-6pm F, 9am-5pm Sat, 1pm-5pm Sun", "10am-5pm daily",
// "12pm-5pm, Mon-Sat" or "12am-7am, 9am-5:30pm, 7:30pm-12am daily". The string is considered a
// match as long as at least one group parses successfully.
func parseMultiDaySegments(input string) ([]OpeningHours, bool, error) {
	var ohs []OpeningHours
	for _, group := range strings.Split(input, ";") {
		groupOhs, err := parseDaySegmentGroup(group)
		if err != nil {
			return nil, true, err
		}

		ohs = append(ohs, groupOhs...)
	}

	if len(ohs) == 0 {
		return nil, false, nil
	}

	return ohs, true, nil
}

var dailyRangeRegexp = regexp.MustCompile(`(?i)^(\d{1,2})(?::([0-5][0-9]))?(?:\s*(am|pm))?(?:\s*-\s*|\s+to\s+)(\d{1,2})(?::([0-5][0-9]))?\s*(am|pm)$`)

// parseDailyRange matches a single 12-hour time range applied to every day of the week, e.g.
// "9am-9pm", "9:30am - 5:00pm" or "5am to 10pm", producing seven OpeningHours, one per weekday.
func parseDailyRange(input string) ([]OpeningHours, bool, error) {
	matches := dailyRangeRegexp.FindStringSubmatch(input)
	if matches == nil {
		return nil, false, nil
	}

	ohs, err := buildEveryDayOpeningHours(matches)
	if err != nil {
		return nil, true, err
	}

	return ohs, true, nil
}

// parse12HourTime parses a 12-hour clock time, given as separate hour, minute and meridiem
// ("am"/"pm") strings, into minutes since midnight. minuteStr may be empty, in which case the
// minutes are assumed to be zero. meridiem may also be empty (eg. an opening time written as
// "9-5pm", where only the closing time carries "pm"), in which case it's assumed to be "am".
func parse12HourTime(hourStr, minuteStr, meridiem string) (int, error) {
	hour, err := strconv.Atoi(hourStr)
	if err != nil || hour < 1 || hour > 12 {
		return 0, fmt.Errorf("invalid hour value")
	}

	minute := 0
	if minuteStr != "" {
		minute, err = strconv.Atoi(minuteStr)
		if err != nil || minute < 0 || minute > 59 {
			return 0, fmt.Errorf("invalid minute value")
		}
	}

	if meridiem == "" {
		meridiem = "am"
	}

	switch strings.ToLower(meridiem) {
	case "am":
		if hour == 12 {
			hour = 0
		}
	case "pm":
		if hour != 12 {
			hour += 12
		}
	default:
		return 0, fmt.Errorf("invalid meridiem value")
	}

	return hour*60 + minute, nil
}

// parse12HourClosingTime parses a 12-hour clock closing time the same way as parse12HourTime,
// except that "12am" is treated as the end of the day (24:00, ie. 1440 minutes since midnight)
// rather than its start. As an opening time, "12am" means midnight at the start of the day (eg.
// "12am-8am" opens at 00:00); as a closing time in the same day/day range, it conventionally means
// open until the end of that day (eg. "10am-12am" closes at 24:00, not 00:00).
func parse12HourClosingTime(hourStr, minuteStr, meridiem string) (int, error) {
	minutes, err := parse12HourTime(hourStr, minuteStr, meridiem)
	if err != nil {
		return 0, err
	}

	if minutes == 0 && strings.EqualFold(meridiem, "am") {
		return 1440, nil
	}

	return minutes, nil
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
