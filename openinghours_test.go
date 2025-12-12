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
		{"sunday", 0, ""},
		{"Monday", 1, ""},
		{"TUESDAY", 2, ""},
		{"friDAY", 5, ""},
		{"", 0, "invalid weekday"},
		{"funday", 0, "invalid weekday"},
		{"mon", 0, "invalid weekday"},
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
