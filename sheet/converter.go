package sheet

import "time"

type Converter struct {
	layout           string
	baseDateString   string
	baseColumnNumber int64
	WeekdayMap       map[string]string
}

func NewConverter() *Converter {
	return &Converter{
		layout:           "2006-01-02",
		baseDateString:   "2020-09-28",
		baseColumnNumber: 129,
		WeekdayMap: map[string]string{
			"Monday":    "月",
			"Tuesday":   "火",
			"Wednesday": "水",
			"Thursday":  "木",
			"Friday":    "金",
			"Saturday":  "土",
			"Sunday":    "日",
		},
	}
}

func (c *Converter) GetColumnNumber(stringDate string) int64 {
	baseDate, _ := time.Parse(c.layout, c.baseDateString)
	date, _ := time.Parse(c.layout, stringDate)
	durationDate := date.Sub(baseDate).Hours() / 24

	return c.baseColumnNumber + int64(durationDate)
}

func (c *Converter) GetDate(num int64) time.Time {
	duration := num - c.baseColumnNumber
	baseDate, _ := time.Parse(c.layout, c.baseDateString)
	date := baseDate.AddDate(0, 0, int(duration))

	return date
}
