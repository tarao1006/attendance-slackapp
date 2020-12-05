package sheet

import (
	"testing"
)

func TestSplitTimes(t *testing.T) {
	tests := []struct {
		Value    string
		Expected AttendanceTime
	}{
		{
			"\n12:00-",
			AttendanceTime{"", "", "12:00", ""},
		},
		{
			"\n-13:00",
			AttendanceTime{"", "", "", "13:00"},
		},
		{
			"12:00-13:00\n12:00-",
			AttendanceTime{"12:00", "13:00", "12:00", ""},
		},
		{
			"12:00-13:00\n-13:00",
			AttendanceTime{"12:00", "13:00", "", "13:00"},
		},
		{
			"12:00-13:00\n12:00-13:00",
			AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
		},
	}

	for _, tt := range tests {
		result := splitTimes(tt.Value)
		if *result != tt.Expected {
			t.Errorf("expected: %s, result: %s", tt.Expected, result)
		}
	}
}

func TestExtractTime(t *testing.T) {
	tests := []struct {
		Value    string
		Expected AttendanceTime
	}{
		{
			"",
			AttendanceTime{"", "", "", ""},
		},
		{
			"12:00-13:00",
			AttendanceTime{"12:00", "13:00", "", ""},
		},
		{
			"\n(12:00-)",
			AttendanceTime{"", "", "12:00", ""},
		},
		{
			"\n(-13:00)",
			AttendanceTime{"", "", "", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-)",
			AttendanceTime{"12:00", "13:00", "12:00", ""},
		},
		{
			"12:00-13:00\n(-13:00)",
			AttendanceTime{"12:00", "13:00", "", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-13:00)",
			AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
		},
	}

	for _, tt := range tests {
		result := ExtractTime(tt.Value)
		if *result != tt.Expected {
			t.Errorf("expected: %s, result: %s", tt.Expected, result)
		}
	}
}

func TestAttendanceTime_Format(t *testing.T) {
	tests := []struct {
		*AttendanceTime
		Expected string
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

	for _, tt := range tests {
		result := tt.Format()
		if result != tt.Expected {
			t.Errorf("expected: %s, result: %s", tt.Expected, result)
		}
	}
}

func TestAddPlan(t *testing.T) {
	tests := []struct {
		CurrentValue  string
		PlanStartTime string
		PlanEndTime   string
		Expected      AttendanceTime
	}{
		{
			"",
			"12:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "", ""},
		},
		{
			"12:00-13:00",
			"12:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "", ""},
		},
		{
			"\n(12:00-)",
			"12:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "12:00", ""},
		},
		{
			"\n(-13:00)",
			"12:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-)",
			"12:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "12:00", ""},
		},
		{
			"12:00-13:00\n(-13:00)",
			"12:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-13:00)",
			"12:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
		},
	}

	for _, tt := range tests {
		result := AddPlan(tt.CurrentValue, tt.PlanStartTime, tt.PlanEndTime)
		if *result != tt.Expected {
			t.Errorf("expected: %s, result: %s", tt.Expected, result)
		}
	}
}

func TestEnteredTime(t *testing.T) {
	tests := []struct {
		CurrentValue string
		StartTime    string
		Expected     AttendanceTime
	}{
		{
			"",
			"12:00",
			AttendanceTime{"", "", "12:00", ""},
		},
		{
			"12:00-13:00",
			"12:00",
			AttendanceTime{"12:00", "13:00", "12:00", ""},
		},
		{
			"\n(12:00-)",
			"12:00",
			AttendanceTime{"", "", "12:00", ""},
		},
		{
			"\n(-13:00)",
			"12:00",
			AttendanceTime{"", "", "12:00", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-)",
			"12:00",
			AttendanceTime{"12:00", "13:00", "12:00", ""},
		},
		{
			"12:00-13:00\n(-13:00)",
			"12:00",
			AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-13:00)",
			"12:00",
			AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
		},
	}

	for _, tt := range tests {
		result := AddEnteredTine(tt.CurrentValue, tt.StartTime)
		if *result != tt.Expected {
			t.Errorf("expected: %s, result: %s", tt.Expected, result)
		}
	}
}

func TestLeftTime(t *testing.T) {
	tests := []struct {
		CurrentValue string
		EndTime      string
		Expected     AttendanceTime
	}{
		{
			"",
			"13:00",
			AttendanceTime{"", "", "", "13:00"},
		},
		{
			"12:00-13:00",
			"13:00",
			AttendanceTime{"12:00", "13:00", "", "13:00"},
		},
		{
			"\n(12:00-)",
			"13:00",
			AttendanceTime{"", "", "12:00", "13:00"},
		},
		{
			"\n(-13:00)",
			"13:00",
			AttendanceTime{"", "", "", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-)",
			"13:00",
			AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
		},
		{
			"12:00-13:00\n(-13:00)",
			"13:00",
			AttendanceTime{"12:00", "13:00", "", "13:00"},
		},
		{
			"12:00-13:00\n(12:00-13:00)",
			"13:00",
			AttendanceTime{"12:00", "13:00", "12:00", "13:00"},
		},
	}

	for _, tt := range tests {
		result := AddLeftTime(tt.CurrentValue, tt.EndTime)
		if *result != tt.Expected {
			t.Errorf("expected: %s, result: %s", tt.Expected, result)
		}
	}
}
