package sheet

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetService struct {
	config        *oauth2.Config
	token         *oauth2.Token
	converter     *Converter
	spreadsheetID string
}

func NewSpreadsheetService() *SpreadsheetService {
	t, err := time.Parse(time.RFC3339, os.Getenv("EXPIRY"))
	if err != nil {
		log.Fatal(err)
	}
	return &SpreadsheetService{
		config: &oauth2.Config{
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
		},
		token: &oauth2.Token{
			AccessToken:  os.Getenv("ACCESS_TOKEN"),
			TokenType:    os.Getenv("TOKEN_TYPE"),
			RefreshToken: os.Getenv("REFRESH_TOKEN"),
			Expiry:       t,
		},
		converter:     NewConverter(),
		spreadsheetID: os.Getenv("SPREADSHEET_ID"),
	}

}

func (s *SpreadsheetService) service() (*sheets.Service, error) {
	ctx := context.Background()
	tokenSource := s.config.TokenSource(ctx, s.token)
	return sheets.NewService(ctx, option.WithTokenSource(tokenSource))
}

func (s *SpreadsheetService) Add(userID string, date string, startTime string, endTime string) {
	s.preExecute(date)
	updateRange := s.getTargetRange(userID, date)
	currentValue := s.getCurrentValue(updateRange)
	attendanceTime := AddPlan(currentValue, startTime, endTime)

	if _, err := s.update(updateRange, attendanceTime.Format()); err != nil {
		log.Fatal(err)
	}
}

func (s *SpreadsheetService) Enter(userID string) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowTime := time.Now().In(jst)
	date := nowTime.Format("2006-01-02")
	s.preExecute(date)
	updateRange := s.getTargetRange(userID, date)
	currentValue := s.getCurrentValue(updateRange)
	attendanceTime := AddEnteredTine(currentValue, nowTime.Format("15:04"))

	if _, err := s.update(updateRange, attendanceTime.Format()); err != nil {
		log.Fatal(err)
	}
}

func (s *SpreadsheetService) Leave(userID string) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowTime := time.Now().In(jst)
	date := nowTime.Format("2006-01-02")
	s.preExecute(date)
	var updateRange = s.getTargetRange(userID, date)
	currentValue := s.getCurrentValue(updateRange)
	var attendanceTime = AddLeftTime(currentValue, nowTime.Format("15:04"))

	if attendanceTime.StartTime == "" {
		yesterdayRange := s.getTargetRange(userID, nowTime.AddDate(0, 0, -1).Format("2006-01-02"))
		yesterdayValue := s.getCurrentValue(yesterdayRange)
		yesterdayAttendanceTime := ExtractTime(yesterdayValue)
		if yesterdayAttendanceTime.StartTime != "" && yesterdayAttendanceTime.EndTime == "" {
			attendanceTime = AddLeftTime(yesterdayValue, nowTime.Format("15:04"))
			updateRange = yesterdayRange
		}
	}

	if _, err := s.update(updateRange, attendanceTime.Format()); err != nil {
		log.Fatal(err)
	}
}

func (s *SpreadsheetService) preExecute(date string) {
	targetColumnNumber := s.converter.GetColumnNumber(date)
	nowColumnCount, sheetID := s.getSheetInfomation()

	if nowColumnCount < targetColumnNumber {
		duration := targetColumnNumber - nowColumnCount
		if _, err := s.appendColumns(sheetID, duration); err != nil {
			log.Fatal(err)
		}
		if _, err := s.writeDate(nowColumnCount, duration); err != nil {
			log.Fatal(err)
		}
	}
}

func (s *SpreadsheetService) getTargetRange(userID string, date string) string {
	targetColumnNumber := s.converter.GetColumnNumber(date)
	return "シート2!" + ConvertIntToString(targetColumnNumber) + os.Getenv(userID)
}

func (s *SpreadsheetService) getSheetInfomation() (int64, int64) {
	srv, err := s.service()
	if err != nil {
		log.Fatal(err)
	}
	respSheet, err := srv.Spreadsheets.Get(s.spreadsheetID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}
	var nowColumnCount int64
	var sheetID int64
	for _, sheet := range respSheet.Sheets {
		if sheet.Properties.Title == "シート2" {
			sheetID = sheet.Properties.SheetId
			nowColumnCount = sheet.Properties.GridProperties.ColumnCount
		}
	}
	return nowColumnCount, sheetID
}

func (s *SpreadsheetService) appendColumns(sheetID int64, duration int64) (*sheets.BatchUpdateSpreadsheetResponse, error) {
	srv, err := s.service()
	if err != nil {
		log.Fatal(err)
	}
	return srv.Spreadsheets.BatchUpdate(s.spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
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

func (s *SpreadsheetService) writeDate(nowColumnCount int64, duration int64) (*sheets.UpdateValuesResponse, error) {
	srv, err := s.service()
	if err != nil {
		log.Fatal(err)
	}

	leftColumnString := ConvertIntToString(nowColumnCount + 1)
	rightColumnString := ConvertIntToString(nowColumnCount + duration)

	targetRange := "シート2!" + leftColumnString + "2:" + rightColumnString + "3"
	dummy := make([][]interface{}, duration)

	for i := 1; i <= int(duration); i++ {
		t := s.converter.GetDate(int64(i) + nowColumnCount)
		dummy[i-1] = []interface{}{t.Format("1/02"), s.converter.WeekdayMap[t.Format("Monday")]}
	}

	return srv.Spreadsheets.Values.Update(
		s.spreadsheetID,
		targetRange,
		&sheets.ValueRange{
			MajorDimension: "COLUMNS",
			Range:          targetRange,
			Values:         dummy,
		},
	).ValueInputOption("USER_ENTERED").Do()
}

func (s *SpreadsheetService) update(updateRange string, newValue string) (*sheets.UpdateValuesResponse, error) {
	srv, err := s.service()
	if err != nil {
		log.Fatal(err)
	}

	return srv.Spreadsheets.Values.Update(
		s.spreadsheetID,
		updateRange,
		&sheets.ValueRange{
			MajorDimension: "COLUMNS",
			Range:          updateRange,
			Values:         [][]interface{}{{newValue}},
		},
	).ValueInputOption("USER_ENTERED").Do()
}

func (s *SpreadsheetService) getCurrentValue(updateRange string) string {
	srv, err := s.service()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := srv.Spreadsheets.Values.Get(
		s.spreadsheetID,
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
