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

	for _, testCase := range testCases {
		result := testCase.Format()
		if result != testCase.expected {
			t.Errorf("expected: %s, result: %s", testCase.expected, result)
		}
	}
}
