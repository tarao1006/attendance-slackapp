package sheet

import (
	"testing"
)

func TestAttendanceTime_Format(t *testing.T) {

	testCases := []struct {
		*AttendanceTime
		expected string
	}{
		{
			&AttendanceTime{"", "", "", ""},
			"",
		},
		{
			&AttendanceTime{"12:00", "13:00", "", ""},
			"12:00-13:00",
		},
		{
			&AttendanceTime{"", "", "12:00", ""},
			"\n(12:00-)",
		},
		{
			&AttendanceTime{"", "", "", "13:00"},
			"\n(-13:00)",
		},
		{
			&AttendanceTime{"12:00", "13:00", "12:00", ""},
			"12:00-13:00\n(12:00-)",
		},
		{
			&AttendanceTime{"12:00", "13:00", "", "13:00"},
			"12:00-13:00\n(-13:00)",
		},
		{
			&AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
			"12:00-13:00\n(12:00-13:00)",
		},
	}

	for _, tt := range testCases {
		result := tt.Format()
		if result != tt.expected {
			t.Errorf("expected: %s, result: %s", tt.expected, result)
		}
	}
}
