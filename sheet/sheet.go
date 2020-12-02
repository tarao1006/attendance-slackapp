package sheet

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const LAYOUT string = "2006-01-02"
const BASE_DATE_STRING string = "2020-09-28"
const BASE_COLUMN_NUMBER int64 = 1 // 129

var TRANSLATE_WEEKDAY = map[string]string{
	"Monday":    "月",
	"Tuesday":   "火",
	"Wednesday": "水",
	"Thursday":  "木",
	"Friday":    "金",
	"Saturday":  "土",
	"Sunday":    "日",
}

func getColumnNumber(date string) int64 {
	baseDate, _ := time.Parse(LAYOUT, BASE_DATE_STRING)
	d, _ := time.Parse(LAYOUT, date)
	durationDate := d.Sub(baseDate).Hours() / 24

	return BASE_COLUMN_NUMBER + int64(durationDate)
}

func getDateFromNumber(num int64) time.Time {
	duration := num - BASE_COLUMN_NUMBER
	baseDate, _ := time.Parse(LAYOUT, BASE_DATE_STRING)
	date := baseDate.AddDate(0, 0, int(duration))

	return date
}

func convertIntToByte(n int64, now []byte) []byte {
	a := (n - 1) / 26
	b := (n - 1) % 26
	if a == 0 {
		return append([]byte{byte(b + 64 + 1)}, now...)
	}
	return convertIntToByte(a, append([]byte{byte(b + 64 + 1)}, now...))
}

func convertIntToString(n int64) string {
	return string(convertIntToByte(n, make([]byte, 0)))
}

func initService() (*sheets.Service, error) {
	ctx := context.Background()
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:   os.Getenv("AUTH_URI"),
			TokenURL:  os.Getenv("TOKEN_URI"),
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		RedirectURL: os.Getenv("REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	t, _ := time.Parse(time.RFC3339, os.Getenv("EXPIRY"))

	tok := &oauth2.Token{
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		TokenType:    os.Getenv("TOKEN_TYPE"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
		Expiry:       t,
	}

	tokenSource := config.TokenSource(ctx, tok)
	return sheets.NewService(ctx, option.WithTokenSource(tokenSource))
}

func getSheetInfomation(srv *sheets.Service, spreadsheetID string) (int64, int64) {
	respSheet, err := srv.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	var nowColumnCount int64
	var sheetID int64
	for _, sheet := range respSheet.Sheets {
		if sheet.Properties.Title == "シート1" {
			sheetID = sheet.Properties.SheetId
			nowColumnCount = sheet.Properties.GridProperties.ColumnCount
		}
	}
	return nowColumnCount, sheetID
}

func appendColumns(srv *sheets.Service, spreadsheetID string, sheetID int64, duration int64) (*sheets.BatchUpdateSpreadsheetResponse, error) {
	return srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			&sheets.Request{
				AppendDimension: &sheets.AppendDimensionRequest{
					SheetId:   sheetID,
					Dimension: "COLUMNS",
					Length:    duration,
				},
			},
		},
	}).Do()
}

func writeDate(srv *sheets.Service, spreadsheetID string, nowColumnCount int64, duration int64) (*sheets.UpdateValuesResponse, error) {
	leftColumnString := convertIntToString(nowColumnCount + 1)
	rightColumnString := convertIntToString(nowColumnCount + duration)

	targetRange := leftColumnString + "2:" + rightColumnString + "3"
	dummy := make([][]interface{}, duration)

	for i := 1; i <= int(duration); i++ {
		t := getDateFromNumber(int64(i) + nowColumnCount)
		dummy[i-1] = []interface{}{t.Format("1/02"), TRANSLATE_WEEKDAY[t.Format("Monday")]}
	}

	return srv.Spreadsheets.Values.Update(
		spreadsheetID,
		targetRange,
		&sheets.ValueRange{
			MajorDimension: "COLUMNS",
			Range:          targetRange,
			Values:         dummy,
		},
	).ValueInputOption("USER_ENTERED").Do()
}

func getNowValue(srv *sheets.Service, spreadsheetID string, updateRange string) string {
	resp, err := srv.Spreadsheets.Values.Get(
		spreadsheetID,
		updateRange,
	).Do()

	if err != nil {
		log.Println(err.Error())
	}

	if len(resp.Values) == 0 {
		return ""
	} else {
		if value, ok := resp.Values[0][0].(string); ok {
			return value
		} else {
			return ""
		}
	}
}

func appendPlan(srv *sheets.Service, spreadsheetID string, updateRange string, startTime string, endTime string) (*sheets.UpdateValuesResponse, error) {
	return srv.Spreadsheets.Values.Update(
		spreadsheetID,
		updateRange,
		&sheets.ValueRange{
			MajorDimension: "COLUMNS",
			Range:          updateRange,
			Values:         [][]interface{}{{startTime + "-" + endTime}},
		},
	).ValueInputOption("USER_ENTERED").Do()
}

func appendEnter(srv *sheets.Service, spreadsheetID string, updateRange string, now string) (*sheets.UpdateValuesResponse, error) {
	return srv.Spreadsheets.Values.Update(
		spreadsheetID,
		updateRange,
		&sheets.ValueRange{
			MajorDimension: "COLUMNS",
			Range:          updateRange,
			Values:         [][]interface{}{{now + "\n(" + time.Now().Format("15:04") + "-)"}},
		},
	).ValueInputOption("USER_ENTERED").Do()
}

func appendLeave(srv *sheets.Service, spreadsheetID string, updateRange string, now string) (*sheets.UpdateValuesResponse, error) {
	values := strings.Split(now, "\n")
	up := values[0]
	low := values[1]
	newLow := up + "\n" + low[:len(low)-1] + time.Now().Format("15:04") + ")"
	return srv.Spreadsheets.Values.Update(
		spreadsheetID,
		updateRange,
		&sheets.ValueRange{
			MajorDimension: "COLUMNS",
			Range:          updateRange,
			Values:         [][]interface{}{{newLow}},
		},
	).ValueInputOption("USER_ENTERED").Do()
}

func Edit(userID string, date string, startTime string, endTime string, operationType string) {
	var targetColumnNumber = getColumnNumber(date)
	var targetColumnString = convertIntToString(targetColumnNumber)
	var targetRowNumber = os.Getenv(userID)
	var updateRange = targetColumnString + targetRowNumber

	srv, err := initService()
	if err != nil {
		log.Fatal(err)
	}

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	nowColumnCount, sheetID := getSheetInfomation(srv, spreadsheetID)

	if nowColumnCount < targetColumnNumber {
		duration := targetColumnNumber - nowColumnCount
		if _, err := appendColumns(srv, spreadsheetID, sheetID, duration); err != nil {
			log.Fatal(err)
		}
		if _, err := writeDate(srv, spreadsheetID, nowColumnCount, duration); err != nil {
			log.Fatal(err)
		}
	}

	nowValue := getNowValue(srv, spreadsheetID, updateRange)

	switch operationType {
	case "add":
		if _, err := appendPlan(srv, spreadsheetID, updateRange, startTime, endTime); err != nil {
			log.Fatal(err)
		}
	case "enter":
		if _, err := appendEnter(srv, spreadsheetID, updateRange, nowValue); err != nil {
			log.Fatal(err)
		}
	case "leave":
		if _, err := appendLeave(srv, spreadsheetID, updateRange, nowValue); err != nil {
			log.Fatal(err)
		}
	}
}
