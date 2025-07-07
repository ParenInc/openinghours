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
		expectedResult map[string][]string
		expectedError  error
	}{
		"Open all Week Long": {
			openingHours: "W1T00:00:00/W7T24:00:00",
			expectedResult: map[string][]string{
				"monday":    []string{"open: 00:00, close: 24:00"},
				"tuesday":   []string{"open: 00:00, close: 24:00"},
				"wednesday": []string{"open: 00:00, close: 24:00"},
				"thursday":  []string{"open: 00:00, close: 24:00"},
				"friday":    []string{"open: 00:00, close: 24:00"},
				"saturday":  []string{"open: 00:00, close: 24:00"},
				"sunday":    []string{"open: 00:00, close: 24:00"},
			},
			expectedError: nil,
		},
		"Open twice on monday": {
			openingHours: "W1T08:00:00/W1T12:00:00,W1T13:00:00/W1T18:00:00",
			expectedResult: map[string][]string{
				"monday": []string{
					"open: 08:00, close: 12:00",
					"open: 13:00, close: 18:00",
				},
			},
			expectedError: nil,
		},
		"never open": {
			openingHours:   "",
			expectedResult: map[string][]string{},
			expectedError:  nil,
		},
		"when string invalid": {
			openingHours:   "invalid",
			expectedResult: nil,
			expectedError:  fmt.Errorf("invalid opening hours string `invalid`"),
		},
		"starts on monday and end on tuesday": {
			openingHours: "W1T08:00:00/W2T16:00:00",
			expectedResult: map[string][]string{
				"monday": []string{
					"open: 08:00, close: 24:00",
				},
				"tuesday": []string{
					"open: 00:00, close: 16:00",
				},
			},
			expectedError: nil,
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := GetHumanReadableTimes(tt.openingHours)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
