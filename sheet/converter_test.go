package sheet

import (
	"testing"
	"time"
)

func TestConverter_GetColumnNumber(t *testing.T) {
	converter := NewConverter()
	tests := []struct {
		Date     string
		Expected int64
	}{
		{
			"2020-09-28",
			129,
		},
		{
			"2020-12-30",
			222,
		},
		{
			"2021-01-01",
			224,
		},
	}

	for _, tt := range tests {
		result := converter.GetColumnNumber(tt.Date)
		if result != tt.Expected {
			t.Errorf("expected: %d, result: %d", tt.Expected, result)
		}
	}
}

func TestConverter_GetDate(t *testing.T) {
	converter := NewConverter()
	tests := []struct {
		Num      int64
		Expected time.Time
	}{
		{
			129,
			time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			222,
			time.Date(2020, time.December, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			224,
			time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		result := converter.GetDate(tt.Num)
		if result != tt.Expected {
			t.Errorf("expected: %v, result: %v", tt.Expected, result)
		}
	}
}

func TestConvertIntToString(t *testing.T) {
	tests := []struct {
		ColumnNumber int64
		Expected     string
	}{
		{
			1,
			"A",
		},
		{
			2,
			"B",
		},
		{
			222,
			"HN",
		},
		{
			666,
			"YP",
		},
		{
			752,
			"ABX",
		},
	}

	for _, tt := range tests {
		result := ConvertIntToString(tt.ColumnNumber)
		if result != tt.Expected {
			t.Errorf("expected: %s, result: %s", tt.Expected, result)
		}
	}
}
