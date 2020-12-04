package sheet

import (
	"fmt"
	"regexp"
	"strings"
)

type AttendanceTime struct {
	PlanStartTime string
	PlanEndTime   string
	StartTime     string
	EndTime       string
}

func (t *AttendanceTime) SprintRow() string {
	return fmt.Sprintf("%s-%s", t.PlanStartTime, t.PlanEndTime)
}

func (t *AttendanceTime) SprintTwoRows() string {
	if t.PlanStartTime == "" && t.PlanEndTime == "" {
		return fmt.Sprintf("\n(%s-%s)", t.StartTime, t.EndTime)
	} else {
		return fmt.Sprintf("%s-%s\n(%s-%s)", t.PlanStartTime, t.PlanEndTime, t.StartTime, t.EndTime)
	}
}

func SplitTimes(value string) *AttendanceTime {
	rows := strings.Split(value, "\n")
	planTimes := strings.Split(rows[0], "-")
	times := strings.Split(rows[1], "-")
	if len(planTimes) == 1 {
		return &AttendanceTime{
			PlanStartTime: "",
			PlanEndTime:   "",
			StartTime:     times[0],
			EndTime:       times[1],
		}
	} else {
		return &AttendanceTime{
			PlanStartTime: planTimes[0],
			PlanEndTime:   planTimes[1],
			StartTime:     times[0],
			EndTime:       times[1],
		}
	}
}

func ExtractTime(value string) *AttendanceTime {
	const timeRegex = "([01][0-9]|2[0-3]):[0-5][0-9]"
	const towTimeRegex = timeRegex + "-" + timeRegex
	patterns := []*regexp.Regexp{
		regexp.MustCompile(``),
		regexp.MustCompile("^" + towTimeRegex + "$"),
		regexp.MustCompile(`^\n\(` + timeRegex + `-\)$`),
		regexp.MustCompile(`^\n\(-` + timeRegex + `\)$`),
		regexp.MustCompile("^" + towTimeRegex + `\n\(` + timeRegex + `-\)$`),
		regexp.MustCompile("^" + towTimeRegex + `\n\(-` + timeRegex + `\)$`),
		regexp.MustCompile("^" + towTimeRegex + `\n\(` + towTimeRegex + `\)$`),
	}

	for i, pattern := range patterns {
		var matchedValue = string(pattern.Find([]byte(value)))

		if matchedValue != "" {
			matchedValue = strings.ReplaceAll(matchedValue, "(", "")
			matchedValue = strings.ReplaceAll(matchedValue, ")", "")
			switch i {
			case 0:
				return &AttendanceTime{}
			case 1:
				times := strings.Split(matchedValue, "-")
				return &AttendanceTime{
					PlanStartTime: times[0],
					PlanEndTime:   times[1],
				}
			case 2, 3, 4, 5, 6:
				return SplitTimes(matchedValue)
			}
		}
	}
	return &AttendanceTime{}
}

func AddPlan(currentValue string, planStartTime string, planEndTime string) string {
	attendanceTime := ExtractTime(currentValue)
	attendanceTime.PlanStartTime = planStartTime
	attendanceTime.PlanEndTime = planEndTime

	if attendanceTime.StartTime == "" && attendanceTime.EndTime == "" {
		return attendanceTime.SprintRow()
	} else {
		return attendanceTime.SprintTwoRows()
	}
}

func AddEnteredTine(currentValue string, startTime string) string {
	attendanceTime := ExtractTime(currentValue)
	attendanceTime.StartTime = startTime
	return attendanceTime.SprintTwoRows()
}

func AddLeftTime(currentValue string, endTime string) string {
	attendanceTime := ExtractTime(currentValue)
	attendanceTime.EndTime = endTime
	return attendanceTime.SprintTwoRows()
}
