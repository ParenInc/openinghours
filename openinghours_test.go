package openinghours

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpeningHoursString(t *testing.T) {
	tests := map[string]struct {
		openingHours   OpeningHours
		expectedResult string
	}{
		"when monday": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 480},
				Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
			},
			expectedResult: "W1T08:00:00/W1T16:00:00",
		},
		"when tuesday": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 360},
				Close: &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 1200},
			},
			expectedResult: "W2T06:00:00/W2T20:00:00",
		},
		"when wednesday": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 3, MinutesSinceMidnight: 480},
				Close: &TimeInWeek{Weekday: 3, MinutesSinceMidnight: 960},
			},
			expectedResult: "W3T08:00:00/W3T16:00:00",
		},
		"when thursday": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 4, MinutesSinceMidnight: 490},
				Close: &TimeInWeek{Weekday: 4, MinutesSinceMidnight: 975},
			},
			expectedResult: "W4T08:10:00/W4T16:15:00",
		},
		"when friday": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 5, MinutesSinceMidnight: 630},
				Close: &TimeInWeek{Weekday: 5, MinutesSinceMidnight: 780},
			},
			expectedResult: "W5T10:30:00/W5T13:00:00",
		},
		"when saturday": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 6, MinutesSinceMidnight: 480},
				Close: &TimeInWeek{Weekday: 6, MinutesSinceMidnight: 960},
			},
			expectedResult: "W6T08:00:00/W6T16:00:00",
		},
		"when sunday": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 7, MinutesSinceMidnight: 480},
				Close: &TimeInWeek{Weekday: 7, MinutesSinceMidnight: 960},
			},
			expectedResult: "W7T08:00:00/W7T16:00:00",
		},
		"when closing time is during the next day": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 1200},
				Close: &TimeInWeek{Weekday: 3, MinutesSinceMidnight: 240},
			},
			expectedResult: "W2T20:00:00/W3T04:00:00",
		},
		"when opening hours not specified": {
			openingHours: OpeningHours{
				Open:  nil,
				Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
			},
			expectedResult: "/W1T16:00:00",
		},
		"when closing hours not specified": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 480},
				Close: nil,
			},
			expectedResult: "W1T08:00:00/",
		},
		"when opening and closing hours not specified": {
			openingHours: OpeningHours{
				Open:  nil,
				Close: nil,
			},
			expectedResult: "/",
		},
		"when weekday invalid": {
			openingHours: OpeningHours{
				Open:  &TimeInWeek{Weekday: 10, MinutesSinceMidnight: 480},
				Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
			},
			expectedResult: "W10T08:00:00/W1T16:00:00",
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := tt.openingHours.String()
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestOpeningHoursSliceToString(t *testing.T) {
	tests := map[string]struct {
		openingHours   []OpeningHours
		expectedResult string
	}{
		"single monday": {
			openingHours: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 480},
					Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
				},
			},
			expectedResult: "W1T08:00:00/W1T16:00:00",
		},
		"single tuesday": {
			openingHours: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 360},
					Close: &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 1200},
				},
			},
			expectedResult: "W2T06:00:00/W2T20:00:00",
		},
		"multiple days": {
			openingHours: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 480},
					Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
				},
				{
					Open:  &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 360},
					Close: &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 1200},
				},
			},
			expectedResult: "W1T08:00:00/W1T16:00:00,W2T06:00:00/W2T20:00:00",
		},
		"empty slice": {
			openingHours:   []OpeningHours{},
			expectedResult: "",
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := OpeningHoursSliceToString(tt.openingHours)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestParseOpeningHours(t *testing.T) {
	tests := map[string]struct {
		openingHours   string
		expectedResult []OpeningHours
		expectedError  error
	}{
		"when monday": {
			openingHours: "W1T08:00:00/W1T16:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 480},
					Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
				},
			},
			expectedError: nil,
		},
		"when tuesday": {
			openingHours: "W2T06:00:00/W2T20:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 360},
					Close: &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 1200},
				},
			},
			expectedError: nil,
		},
		"when wednesday": {
			openingHours: "W3T08:00:00/W3T16:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 3, MinutesSinceMidnight: 480},
					Close: &TimeInWeek{Weekday: 3, MinutesSinceMidnight: 960},
				},
			},
			expectedError: nil,
		},
		"when thursday": {
			openingHours: "W4T08:10:00/W4T16:15:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 4, MinutesSinceMidnight: 490},
					Close: &TimeInWeek{Weekday: 4, MinutesSinceMidnight: 975},
				},
			},
			expectedError: nil,
		},
		"when friday": {
			openingHours: "W5T10:30:00/W5T13:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 5, MinutesSinceMidnight: 630},
					Close: &TimeInWeek{Weekday: 5, MinutesSinceMidnight: 780},
				},
			},
			expectedError: nil,
		},
		"when saturday": {
			openingHours: "W6T08:00:00/W6T16:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 6, MinutesSinceMidnight: 480},
					Close: &TimeInWeek{Weekday: 6, MinutesSinceMidnight: 960},
				},
			},
			expectedError: nil,
		},
		"when sunday": {
			openingHours: "W7T08:00:00/W7T16:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 7, MinutesSinceMidnight: 480},
					Close: &TimeInWeek{Weekday: 7, MinutesSinceMidnight: 960},
				},
			},
			expectedError: nil,
		},
		"when closing time is during the next day": {
			openingHours: "W2T20:00:00/W3T04:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 1200},
					Close: &TimeInWeek{Weekday: 3, MinutesSinceMidnight: 240},
				},
			},
			expectedError: nil,
		},
		"when opening hours not specified": {
			openingHours: "/W1T16:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  nil,
					Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
				},
			},
			expectedError: nil,
		},
		"when closing hours not specified": {
			openingHours: "W1T08:00:00/",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 480},
					Close: nil,
				},
			},
			expectedError: nil,
		},
		"when opening and closing hours not specified": {
			openingHours: "/",
			expectedResult: []OpeningHours{
				{
					Open:  nil,
					Close: nil,
				},
			},
			expectedError: nil,
		},
		"when whole week": {
			openingHours: "W1T00:00:00/W7T24:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 0},
					Close: &TimeInWeek{Weekday: 7, MinutesSinceMidnight: 1440},
				},
			},
			expectedError: nil,
		},
		"when multiple opening hours": {
			openingHours: "W1T08:00:00/W1T16:00:00,W2T06:00:00/W2T20:00:00",
			expectedResult: []OpeningHours{
				{
					Open:  &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 480},
					Close: &TimeInWeek{Weekday: 1, MinutesSinceMidnight: 960},
				},
				{
					Open:  &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 360},
					Close: &TimeInWeek{Weekday: 2, MinutesSinceMidnight: 1200},
				},
			},
			expectedError: nil,
		},
		"when string empty": {
			openingHours:   "",
			expectedResult: []OpeningHours{},
			expectedError:  nil,
		},
		"when string invalid": {
			openingHours:   "invalid",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid opening hours string `invalid`"),
		},
		"when opening string invalid": {
			openingHours:   "invalid/W1T16:00:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid opening hours: invalid value `invalid`"),
		},
		"when opening weekday invalid": {
			openingHours:   "W9T08:00:00/W1T16:00:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid opening hours: invalid workday in `W9T08:00:00`: expected to be between 1 (monday) and 7 (sunday)"),
		},
		"when opening hours invalid": {
			openingHours:   "W1T99:00:00/W1T16:00:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid opening hours: invalid time in `W1T99:00:00`: invalid hours value"),
		},
		"when opening minutes invalid": {
			openingHours:   "W1T08:99:00/W1T16:00:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid opening hours: invalid time in `W1T08:99:00`: invalid minutes value"),
		},
		"when closing string invalid": {
			openingHours:   "W1T08:00:00/invalid",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid closing hours: invalid value `invalid`"),
		},
		"when closing weekday invalid": {
			openingHours:   "W1T08:00:00/W9T16:00:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid closing hours: invalid workday in `W9T16:00:00`: expected to be between 1 (monday) and 7 (sunday)"),
		},
		"when closing hours invalid": {
			openingHours:   "W1T08:00:00/W1T99:00:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid closing hours: invalid time in `W1T99:00:00`: invalid hours value"),
		},
		"when closing minutes invalid": {
			openingHours:   "W1T08:00:00/W1T16:99:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid closing hours: invalid time in `W1T16:99:00`: invalid minutes value"),
		},
		"when closing time invalid": {
			openingHours:   "W1T00:00:00/W7T24:01:00",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid closing hours: invalid time in `W7T24:01:00`: invalid value"),
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ParseOpeningHours(tt.openingHours)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetHumanReadableTimes(t *testing.T) {
	tests := map[string]struct {
		openingHours   string
		expectedResult map[string][]TimeRange
		expectedError  error
	}{
		"Open all Week Long": {
			openingHours: "W1T00:00:00/W7T24:00:00",
			expectedResult: map[string][]TimeRange{
				"monday":    {{Open: "00:00", Close: "24:00"}},
				"tuesday":   {{Open: "00:00", Close: "24:00"}},
				"wednesday": {{Open: "00:00", Close: "24:00"}},
				"thursday":  {{Open: "00:00", Close: "24:00"}},
				"friday":    {{Open: "00:00", Close: "24:00"}},
				"saturday":  {{Open: "00:00", Close: "24:00"}},
				"sunday":    {{Open: "00:00", Close: "24:00"}},
			},
			expectedError: nil,
		},
		"Open twice on monday": {
			openingHours: "W1T08:00:00/W1T12:00:00,W1T13:00:00/W1T18:00:00",
			expectedResult: map[string][]TimeRange{
				"monday": {
					{Open: "08:00", Close: "12:00"},
					{Open: "13:00", Close: "18:00"},
				},
			},
			expectedError: nil,
		},
		"never open": {
			openingHours:   "",
			expectedResult: nil,
			expectedError:  nil,
		},
		"starts on monday and end on tuesday": {
			openingHours: "W1T08:00:00/W2T16:00:00",
			expectedResult: map[string][]TimeRange{
				"monday": {{
					Open: "08:00", Close: "24:00",
				}},
				"tuesday": {{
					Open: "00:00", Close: "16:00",
				}},
			},
			expectedError: nil,
		},
		"starts on sunday and end on monday at 00:00": {
			openingHours: "W7T00:00:00/W1T00:00:00",
			expectedResult: map[string][]TimeRange{
				"sunday": {{Open: "00:00", Close: "24:00"}},
			},
		},
		"starts on sunday and end on monday": {
			openingHours: "W7T00:00:00/W1T10:00:00",
			expectedResult: map[string][]TimeRange{
				"sunday": {
					{Open: "00:00", Close: "24:00"},
				},
				"monday": {
					{Open: "00:00", Close: "10:00"},
				},
			},
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ohs, err := ParseOpeningHours(tt.openingHours)
			assert.Equal(t, tt.expectedError, err)
			result := GetHumanReadableTimes(ohs)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetOCPIOpeningTimes(t *testing.T) {
	tests := map[string]struct {
		openingHours   string
		expectedResult OCPIOpeningTimes
	}{
		"when 24/7": {
			openingHours: "W1T00:00:00/W7T24:00:00",
			expectedResult: OCPIOpeningTimes{
				TwentyFourSeven: true,
			},
		},
		"when same opening times monday to friday": {
			openingHours: "W1T08:00:00/W1T16:00:00,W2T08:00:00/W2T16:00:00,W3T08:00:00/W3T16:00:00,W4T08:00:00/W4T16:00:00,W5T08:00:00/W5T16:00:00",
			expectedResult: OCPIOpeningTimes{
				TwentyFourSeven: false,
				RegularHours: &[]OCPIRegularHours{
					{
						Weekday:     1,
						PeriodBegin: "08:00",
						PeriodEnd:   "16:00",
					},
					{
						Weekday:     2,
						PeriodBegin: "08:00",
						PeriodEnd:   "16:00",
					},
					{
						Weekday:     3,
						PeriodBegin: "08:00",
						PeriodEnd:   "16:00",
					},
					{
						Weekday:     4,
						PeriodBegin: "08:00",
						PeriodEnd:   "16:00",
					},
					{
						Weekday:     5,
						PeriodBegin: "08:00",
						PeriodEnd:   "16:00",
					},
				},
			},
		},
		"when starts on monday and ends on tuesday": {
			openingHours: "W1T08:00:00/W2T16:00:00",
			expectedResult: OCPIOpeningTimes{
				TwentyFourSeven: false,
				RegularHours: &[]OCPIRegularHours{
					{
						Weekday:     1,
						PeriodBegin: "08:00",
						PeriodEnd:   "00:00",
					},
					{
						Weekday:     2,
						PeriodBegin: "00:00",
						PeriodEnd:   "16:00",
					},
				},
			},
		},
		"when starts on sunday and ends on monday at 00:00": {
			openingHours: "W7T00:00:00/W1T00:00:00",
			expectedResult: OCPIOpeningTimes{
				TwentyFourSeven: false,
				RegularHours: &[]OCPIRegularHours{
					{
						Weekday:     7,
						PeriodBegin: "00:00",
						PeriodEnd:   "00:00",
					},
				},
			},
		},
		"when starts on sunday and ends on monday": {
			openingHours: "W7T00:00:00/W1T10:00:00",
			expectedResult: OCPIOpeningTimes{
				TwentyFourSeven: false,
				RegularHours: &[]OCPIRegularHours{
					{
						Weekday:     7,
						PeriodBegin: "00:00",
						PeriodEnd:   "00:00",
					},
					{
						Weekday:     1,
						PeriodBegin: "00:00",
						PeriodEnd:   "10:00",
					},
				},
			},
		},
		"when opening hours are empty": {
			openingHours:   "",
			expectedResult: OCPIOpeningTimes{},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ohs, _ := ParseOpeningHours(tt.openingHours)
			result := GetOCPIOpeningTimes(ohs)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestParseStringWeekdayToTimeWeekday(t *testing.T) {
	tests := []struct {
		input         string
		expected      int
		expectedError string
	}{
		{"monday", 1, ""},
		{"tuesday", 2, ""},
		{"wednesday", 3, ""},
		{"thursday", 4, ""},
		{"friday", 5, ""},
		{"saturday", 6, ""},
		{"sunday", 7, ""},
		{"Monday", 1, ""},
		{"TUESDAY", 2, ""},
		{"fri", 5, ""},
		{"", 0, "invalid weekday"},
		{"funday", 0, "invalid weekday"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			result, err := ParseStringWeekdayToTimeWeekday(tt.input)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, int(result))
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestParseFreeformOpeningHours(t *testing.T) {
	dailyOpeningHours := func(open, close int) []OpeningHours {
		ohs := make([]OpeningHours, 0, 7)
		for weekday := 1; weekday <= 7; weekday++ {
			ohs = append(ohs, OpeningHours{
				Open:  &TimeInWeek{Weekday: weekday, MinutesSinceMidnight: open},
				Close: &TimeInWeek{Weekday: weekday, MinutesSinceMidnight: close},
			})
		}
		return ohs
	}

	dayOpeningHours := func(weekday, open, close int) OpeningHours {
		return OpeningHours{
			Open:  &TimeInWeek{Weekday: weekday, MinutesSinceMidnight: open},
			Close: &TimeInWeek{Weekday: weekday, MinutesSinceMidnight: close},
		}
	}

	overnightDayOpeningHours := func(weekday, open, close int) OpeningHours {
		closeWeekday := weekday%7 + 1
		return OpeningHours{
			Open:  &TimeInWeek{Weekday: weekday, MinutesSinceMidnight: open},
			Close: &TimeInWeek{Weekday: closeWeekday, MinutesSinceMidnight: close},
		}
	}

	overnightDailyOpeningHours := func(open, close int) []OpeningHours {
		ohs := make([]OpeningHours, 0, 7)
		for weekday := 1; weekday <= 7; weekday++ {
			ohs = append(ohs, overnightDayOpeningHours(weekday, open, close))
		}
		return ohs
	}

	tests := map[string]struct {
		input          string
		expectedResult []OpeningHours
		expectedError  string
	}{
		"9am-9pm": {
			input:          "9am-9pm",
			expectedResult: dailyOpeningHours(540, 1260),
		},
		"24 hours daily": {
			input:          "24 hours daily",
			expectedResult: []OpeningHours{TwentyFourSevenOH},
		},
		"24 Hours Daily mixed case": {
			input:          "24 Hours Daily",
			expectedResult: []OpeningHours{TwentyFourSevenOH},
		},
		"24 hour daily singular": {
			input:          "24 hour daily",
			expectedResult: []OpeningHours{TwentyFourSevenOH},
		},
		"24hoursdaily no spaces": {
			input:          "24hoursdaily",
			expectedResult: []OpeningHours{TwentyFourSevenOH},
		},
		"24/7": {
			input:          "24/7",
			expectedResult: []OpeningHours{TwentyFourSevenOH},
		},
		"24 / 7 with spaces": {
			input:          "24 / 7",
			expectedResult: []OpeningHours{TwentyFourSevenOH},
		},
		"9am-9pm with spaces around dash": {
			input:          "9am - 9pm",
			expectedResult: dailyOpeningHours(540, 1260),
		},
		"9AM-9PM uppercase": {
			input:          "9AM-9PM",
			expectedResult: dailyOpeningHours(540, 1260),
		},
		"9:30am-5:00pm with minutes": {
			input:          "9:30am-5:00pm",
			expectedResult: dailyOpeningHours(570, 1020),
		},
		"12am-12pm midnight to noon": {
			input:          "12am-12pm",
			expectedResult: dailyOpeningHours(0, 720),
		},
		"leading and trailing whitespace": {
			input:          "  9am-9pm  ",
			expectedResult: dailyOpeningHours(540, 1260),
		},
		"empty string": {
			input:         "",
			expectedError: "empty opening hours string",
		},
		"unrecognized format": {
			input:         "open all day",
			expectedError: `unrecognized opening hours format: "open all day"`,
		},
		"missing opening meridiem defaults to am": {
			input:          "9-9pm",
			expectedResult: dailyOpeningHours(540, 1260),
		},
		"closing time before opening time wraps to the next day": {
			input:          "9pm-9am",
			expectedResult: overnightDailyOpeningHours(1260, 540),
		},
		"closing time equal to opening time": {
			input:         "9am-9am",
			expectedError: `closing time must be after opening time in "9am-9am"`,
		},
		"invalid hour": {
			input:         "13am-9pm",
			expectedError: `invalid opening time in "13am-9pm": invalid hour value`,
		},
		"invalid minute": {
			input:         "9:60am-9pm",
			expectedError: `unrecognized opening hours format: "9:60am-9pm"`,
		},
		"12am closing time means end of day, not start of day": {
			input:          "11am-12am",
			expectedResult: dailyOpeningHours(660, 1440),
		},
		"12am opening time still means start of day": {
			input:          "12am-8am",
			expectedResult: dailyOpeningHours(0, 480),
		},
		"12am closing time end of day with a day range": {
			input: "10am-12am M-Sat",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 600, 1440),
				dayOpeningHours(2, 600, 1440),
				dayOpeningHours(3, 600, 1440),
				dayOpeningHours(4, 600, 1440),
				dayOpeningHours(5, 600, 1440),
				dayOpeningHours(6, 600, 1440),
			},
		},
		"daily keyword applies to every weekday": {
			input:          "10am-5pm daily",
			expectedResult: dailyOpeningHours(600, 1020),
		},
		"daily keyword uppercase": {
			input:          "10am-5pm DAILY",
			expectedResult: dailyOpeningHours(600, 1020),
		},
		"daily keyword with trailing notes": {
			input:          "10am-10pm daily; for guest use only; Drivers must bring their own J1772 cordset for Level 1 charging",
			expectedResult: dailyOpeningHours(600, 1320),
		},
		"daily keyword with noon opening time": {
			input:          "12pm-9pm daily",
			expectedResult: dailyOpeningHours(720, 1260),
		},
		"trailing parenthetical aside is ignored": {
			input: "10am-5pm Wed-Sun (May-Sept)",
			expectedResult: []OpeningHours{
				dayOpeningHours(3, 600, 1020),
				dayOpeningHours(4, 600, 1020),
				dayOpeningHours(5, 600, 1020),
				dayOpeningHours(6, 600, 1020),
				dayOpeningHours(7, 600, 1020),
			},
		},
		"parenthetical aside on the last of several clauses is ignored": {
			input: "10am-8pm M-Th, 10am-6pm F, 10am-5pm Sat, 1pm-5pm Sun (Sunday hours in winter only)",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 600, 1200),
				dayOpeningHours(2, 600, 1200),
				dayOpeningHours(3, 600, 1200),
				dayOpeningHours(4, 600, 1200),
				dayOpeningHours(5, 600, 1080),
				dayOpeningHours(6, 600, 1020),
				dayOpeningHours(7, 780, 1020),
			},
		},
		"comma may separate a time range from its day spec": {
			input: "12pm-5pm, Mon-Sat; ask for key card at reception",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 720, 1020),
				dayOpeningHours(2, 720, 1020),
				dayOpeningHours(3, 720, 1020),
				dayOpeningHours(4, 720, 1020),
				dayOpeningHours(5, 720, 1020),
				dayOpeningHours(6, 720, 1020),
			},
		},
		"day spec joining a range and a single day with and": {
			input: "11am-6pm W-F and Sun, 11am-8pm Sat; Winery customers only",
			expectedResult: []OpeningHours{
				dayOpeningHours(3, 660, 1080),
				dayOpeningHours(4, 660, 1080),
				dayOpeningHours(5, 660, 1080),
				dayOpeningHours(7, 660, 1080),
				dayOpeningHours(6, 660, 1200),
			},
		},
		"day spec joining single days with comma and and": {
			input: "9am-5pm Mon, Wed and Fri",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 540, 1020),
				dayOpeningHours(3, 540, 1020),
				dayOpeningHours(5, 540, 1020),
			},
		},
		"unresolvable token within an and-joined day spec is dropped, not an error": {
			input: "9am-5pm Mon and Xyz",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 540, 1020),
			},
		},
		"day spec joining days with ampersand": {
			input: "10am-6pm M & W, 10am-8pm T & Th, 10am-5pm F-Sat",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 600, 1080),
				dayOpeningHours(3, 600, 1080),
				dayOpeningHours(2, 600, 1200),
				dayOpeningHours(4, 600, 1200),
				dayOpeningHours(5, 600, 1020),
				dayOpeningHours(6, 600, 1020),
			},
		},
		"multiple bare time windows share a single trailing daily keyword": {
			input: "12am-7am, 9am-5:30pm, 7:30pm-12am daily",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 0, 420), dayOpeningHours(1, 540, 1050), dayOpeningHours(1, 1170, 1440),
				dayOpeningHours(2, 0, 420), dayOpeningHours(2, 540, 1050), dayOpeningHours(2, 1170, 1440),
				dayOpeningHours(3, 0, 420), dayOpeningHours(3, 540, 1050), dayOpeningHours(3, 1170, 1440),
				dayOpeningHours(4, 0, 420), dayOpeningHours(4, 540, 1050), dayOpeningHours(4, 1170, 1440),
				dayOpeningHours(5, 0, 420), dayOpeningHours(5, 540, 1050), dayOpeningHours(5, 1170, 1440),
				dayOpeningHours(6, 0, 420), dayOpeningHours(6, 540, 1050), dayOpeningHours(6, 1170, 1440),
				dayOpeningHours(7, 0, 420), dayOpeningHours(7, 540, 1050), dayOpeningHours(7, 1170, 1440),
			},
		},
		"multiple bare time windows share a single trailing day range": {
			input: "9am-12pm, 1pm-5pm Sat-Sun",
			expectedResult: []OpeningHours{
				dayOpeningHours(6, 540, 720), dayOpeningHours(6, 780, 1020),
				dayOpeningHours(7, 540, 720), dayOpeningHours(7, 780, 1020),
			},
		},
		"overnight range with an explicit weekday range wraps to the next day": {
			input: "3pm-7:45am M-F",
			expectedResult: []OpeningHours{
				overnightDayOpeningHours(1, 900, 465),
				overnightDayOpeningHours(2, 900, 465),
				overnightDayOpeningHours(3, 900, 465),
				overnightDayOpeningHours(4, 900, 465),
				overnightDayOpeningHours(5, 900, 465),
			},
		},
		"overnight weekday range combined with a 24-hours weekend clause": {
			input: "3pm-7:45am M-F, 24 hours Sat-Sun",
			expectedResult: []OpeningHours{
				overnightDayOpeningHours(1, 900, 465),
				overnightDayOpeningHours(2, 900, 465),
				overnightDayOpeningHours(3, 900, 465),
				overnightDayOpeningHours(4, 900, 465),
				overnightDayOpeningHours(5, 900, 465),
				dayOpeningHours(6, 0, 1440),
				dayOpeningHours(7, 0, 1440),
			},
		},
		"24 hours day spec on its own": {
			input: "24 hours Sat-Sun",
			expectedResult: []OpeningHours{
				dayOpeningHours(6, 0, 1440),
				dayOpeningHours(7, 0, 1440),
			},
		},
		"to used as a synonym for the time-range hyphen": {
			input:          "5am to 10pm daily",
			expectedResult: dailyOpeningHours(300, 1320),
		},
		"to combined with an overnight range and a dropped note": {
			input:          "6am to 2am; pay lot",
			expectedResult: overnightDailyOpeningHours(360, 120),
		},
		"missing opening meridiem with an explicit day defaults to am": {
			input: "9-5pm Sat",
			expectedResult: []OpeningHours{
				dayOpeningHours(6, 540, 1020),
			},
		},
		"missing opening meridiem combined with to and a separate daily clause": {
			input:          "4:30 to 7am; daily",
			expectedResult: dailyOpeningHours(270, 420),
		},
		"ampersand joining two time windows that share a trailing daily keyword": {
			input: "6:30am-6pm & 8pm-10pm daily; see Pro Shop or Restaurant attendant for access",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 390, 1080), dayOpeningHours(1, 1200, 1320),
				dayOpeningHours(2, 390, 1080), dayOpeningHours(2, 1200, 1320),
				dayOpeningHours(3, 390, 1080), dayOpeningHours(3, 1200, 1320),
				dayOpeningHours(4, 390, 1080), dayOpeningHours(4, 1200, 1320),
				dayOpeningHours(5, 390, 1080), dayOpeningHours(5, 1200, 1320),
				dayOpeningHours(6, 390, 1080), dayOpeningHours(6, 1200, 1320),
				dayOpeningHours(7, 390, 1080), dayOpeningHours(7, 1200, 1320),
			},
		},
		"free-text note containing its own comma is dropped instead of erroring": {
			input: "8:30am-5:30pm M-F, 9am-3pm Sat, for client use only",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 510, 1050),
				dayOpeningHours(2, 510, 1050),
				dayOpeningHours(3, 510, 1050),
				dayOpeningHours(4, 510, 1050),
				dayOpeningHours(5, 510, 1050),
				dayOpeningHours(6, 540, 900),
			},
		},
		"free-text note glued via a non-splitting comma after a day range": {
			input: "8am-6pm M-F, for employee use only",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 480, 1080),
				dayOpeningHours(2, 480, 1080),
				dayOpeningHours(3, 480, 1080),
				dayOpeningHours(4, 480, 1080),
				dayOpeningHours(5, 480, 1080),
			},
		},
		"redundant trailing daily keyword after an explicit day range is dropped": {
			input: "9am-9pm M-F daily",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 540, 1260),
				dayOpeningHours(2, 540, 1260),
				dayOpeningHours(3, 540, 1260),
				dayOpeningHours(4, 540, 1260),
				dayOpeningHours(5, 540, 1260),
			},
		},
		"free-text note glued directly onto daily with no separator": {
			input:          "8am-11pm daily during the summer",
			expectedResult: dailyOpeningHours(480, 1380),
		},
		"weekdays keyword": {
			input: "8am-5pm weekdays; tenant use only",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 480, 1020),
				dayOpeningHours(2, 480, 1020),
				dayOpeningHours(3, 480, 1020),
				dayOpeningHours(4, 480, 1020),
				dayOpeningHours(5, 480, 1020),
			},
		},
		"every day keyword": {
			input:          "7am-11:59pm every day; first come, first served for permit holders, university vehicles, and public",
			expectedResult: dailyOpeningHours(420, 1439),
		},
		"closed day dropped when glued onto a preceding day spec": {
			input: "8am-5pm T-Sat, 12pm-5pm Sun, closed M",
			expectedResult: []OpeningHours{
				dayOpeningHours(2, 480, 1020),
				dayOpeningHours(3, 480, 1020),
				dayOpeningHours(4, 480, 1020),
				dayOpeningHours(5, 480, 1020),
				dayOpeningHours(6, 480, 1020),
				dayOpeningHours(7, 720, 1020),
			},
		},
		"S-S means Saturday-Sunday": {
			input: "7am - 10pm S-S, 6am - 10pm M-F",
			expectedResult: []OpeningHours{
				dayOpeningHours(6, 420, 1320),
				dayOpeningHours(7, 420, 1320),
				dayOpeningHours(1, 360, 1320),
				dayOpeningHours(2, 360, 1320),
				dayOpeningHours(3, 360, 1320),
				dayOpeningHours(4, 360, 1320),
				dayOpeningHours(5, 360, 1320),
			},
		},
		"and-joined time windows sharing a trailing day range": {
			input: "8am-12pm and 1pm-5pm M-F; for guest use only",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 480, 720), dayOpeningHours(1, 780, 1020),
				dayOpeningHours(2, 480, 720), dayOpeningHours(2, 780, 1020),
				dayOpeningHours(3, 480, 720), dayOpeningHours(3, 780, 1020),
				dayOpeningHours(4, 480, 720), dayOpeningHours(4, 780, 1020),
				dayOpeningHours(5, 480, 720), dayOpeningHours(5, 780, 1020),
			},
		},
		"bare time range with no day spec at all still applies daily when followed only by a note": {
			input:          "9am-5pm; ask for office key at the front desk",
			expectedResult: dailyOpeningHours(540, 1020),
		},
		"library hours with mixed single and multi-letter day abbreviations and notes": {
			input: "9am-9pm M-Th, 9am-6pm F, 9am-5pm Sat, 1pm-5pm Sun; for library visitors; 2 hour charging limit",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 540, 1260),
				dayOpeningHours(2, 540, 1260),
				dayOpeningHours(3, 540, 1260),
				dayOpeningHours(4, 540, 1260),
				dayOpeningHours(5, 540, 1080),
				dayOpeningHours(6, 540, 1020),
				dayOpeningHours(7, 780, 1020),
			},
		},
		"day ranges mixing single-letter and three-letter abbreviations": {
			input: "9am-9pm M-Th, 9am-6pm F-Sat, 11am-5pm Sun",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 540, 1260),
				dayOpeningHours(2, 540, 1260),
				dayOpeningHours(3, 540, 1260),
				dayOpeningHours(4, 540, 1260),
				dayOpeningHours(5, 540, 1080),
				dayOpeningHours(6, 540, 1080),
				dayOpeningHours(7, 660, 1020),
			},
		},
		"semicolon-separated three-letter day ranges with trailing note": {
			input: "9am-9pm Mon-Tue; 9am-6pm Wed-Sat; for client use only",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 540, 1260),
				dayOpeningHours(2, 540, 1260),
				dayOpeningHours(3, 540, 1080),
				dayOpeningHours(4, 540, 1080),
				dayOpeningHours(5, 540, 1080),
				dayOpeningHours(6, 540, 1080),
			},
		},
		"bare T means Tuesday": {
			input: "9am-9pm M-T",
			expectedResult: []OpeningHours{
				dayOpeningHours(1, 540, 1260),
				dayOpeningHours(2, 540, 1260),
			},
		},
		"bare T range with three-letter abbreviations": {
			input: "10am-5pm T-F, 10am-5pm Sat-Sun",
			expectedResult: []OpeningHours{
				dayOpeningHours(2, 600, 1020),
				dayOpeningHours(3, 600, 1020),
				dayOpeningHours(4, 600, 1020),
				dayOpeningHours(5, 600, 1020),
				dayOpeningHours(6, 600, 1020),
				dayOpeningHours(7, 600, 1020),
			},
		},
		"bare T alongside Th in the same string": {
			input: "10am-5pm T-TH, 10am-9pm F, 10am-5pm Sat-Sun",
			expectedResult: []OpeningHours{
				dayOpeningHours(2, 600, 1020),
				dayOpeningHours(3, 600, 1020),
				dayOpeningHours(4, 600, 1020),
				dayOpeningHours(5, 600, 1260),
				dayOpeningHours(6, 600, 1020),
				dayOpeningHours(7, 600, 1020),
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := ParseFreeformOpeningHours(tt.input)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, result)
			}
		})
	}
}

func TestFreeformOpeningHoursToISOString(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedResult string
		expectedError  string
	}{
		"9am-9pm": {
			input:          "9am-9pm",
			expectedResult: "W1T09:00:00/W1T21:00:00,W2T09:00:00/W2T21:00:00,W3T09:00:00/W3T21:00:00,W4T09:00:00/W4T21:00:00,W5T09:00:00/W5T21:00:00,W6T09:00:00/W6T21:00:00,W7T09:00:00/W7T21:00:00",
		},
		"24 hours daily": {
			input:          "24 hours daily",
			expectedResult: TwentyFourSevenString,
		},
		"invalid input": {
			input:         "nonsense",
			expectedError: `unrecognized opening hours format: "nonsense"`,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := FreeformOpeningHoursToISOString(tt.input)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Equal(t, "", result)
			}
		})
	}
}

func TestParse12HourTime(t *testing.T) {
	tests := []struct {
		hourStr        string
		minuteStr      string
		meridiem       string
		expectedResult int
		expectedError  string
	}{
		{"9", "", "am", 540, ""},
		{"9", "", "pm", 1260, ""},
		{"12", "", "am", 0, ""},
		{"12", "", "pm", 720, ""},
		{"12", "30", "am", 30, ""},
		{"9", "30", "pm", 1290, ""},
		{"0", "", "am", 0, "invalid hour value"},
		{"13", "", "am", 0, "invalid hour value"},
		{"9", "60", "am", 0, "invalid minute value"},
		{"9", "", "xx", 0, "invalid meridiem value"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%s:%s%s", tt.hourStr, tt.minuteStr, tt.meridiem), func(t *testing.T) {
			t.Parallel()
			result, err := parse12HourTime(tt.hourStr, tt.minuteStr, tt.meridiem)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestParse12HourClosingTime(t *testing.T) {
	tests := []struct {
		hourStr        string
		minuteStr      string
		meridiem       string
		expectedResult int
		expectedError  string
	}{
		{"12", "", "am", 1440, ""},
		{"12", "", "pm", 720, ""},
		{"9", "", "pm", 1260, ""},
		{"9", "", "am", 540, ""},
		{"13", "", "am", 0, "invalid hour value"},
		{"9", "60", "am", 0, "invalid minute value"},
		{"9", "", "xx", 0, "invalid meridiem value"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%s:%s%s", tt.hourStr, tt.minuteStr, tt.meridiem), func(t *testing.T) {
			t.Parallel()
			result, err := parse12HourClosingTime(tt.hourStr, tt.minuteStr, tt.meridiem)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestParseMinutesSinceMidnight(t *testing.T) {
	tests := []struct {
		hourStr        string
		minuteStr      string
		expectedResult int
		expectedError  string
	}{
		{"00", "00", 0, ""},
		{"08", "30", 510, ""},
		{"23", "59", 1439, ""},
		{"24", "00", 1440, ""},
		{"24", "01", 0, "invalid value"},
		{"-1", "00", 0, "invalid hours value"},
		{"25", "00", 0, "invalid hours value"},
		{"12", "-1", 0, "invalid minutes value"},
		{"12", "60", 0, "invalid minutes value"},
		{"aa", "00", 0, "invalid hours value"},
		{"12", "bb", 0, "invalid minutes value"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%s:%s", tt.hourStr, tt.minuteStr), func(t *testing.T) {
			t.Parallel()
			result, err := ParseMinutesSinceMidnight(tt.hourStr, tt.minuteStr)
			if tt.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}
