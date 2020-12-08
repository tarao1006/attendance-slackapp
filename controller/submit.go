package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/sheet"
)

type Submit struct {
	client             *slack.Client
	spreadsheetService *sheet.SpreadsheetService
}

func NewSubmit(client *slack.Client, spreadsheetService *sheet.SpreadsheetService) *Submit {
	return &Submit{
		client:             client,
		spreadsheetService: spreadsheetService,
	}
}

func (submit *Submit) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(r.FormValue("payload")), &payload); err != nil {
		log.Printf("Could not parse action response JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := payload.User.ID
	userName := payload.User.Name
	dateString := payload.View.State.Values["date"]["date"].SelectedDate
	startTimeString := payload.View.State.Values["start_time"]["startTime"].Value
	endTimeString := payload.View.State.Values["end_time"]["endTime"].Value

	message := fmt.Sprintf("%s が予定を追加しました\nDate: %s\nStart Time: %s\nEnd Time: %s", userName, dateString, startTimeString, endTimeString)

	date, _ := time.Parse("2006-01-02", dateString)
	if date.Before(time.Now()) {
		resp, _ := json.Marshal(slack.NewErrorsViewSubmissionResponse(map[string]string{
			"date": "過去の日付に予定は追加できません。",
		}))
		w.Header().Add("Content-Type", "application/json")
		w.Write(resp)
		return
	}
	errorMessage := make(map[string]string)
	startTime, err := time.Parse("15:04", startTimeString)
	if err != nil {
		errorMessage["start_time"] = "不正な入力です。"
	}
	endTime, err := time.Parse("15:04", endTimeString)
	if err != nil {
		errorMessage["end_time"] = "不正な入力です。"
	}
	if len(errorMessage) != 0 {
		resp, _ := json.Marshal(slack.NewErrorsViewSubmissionResponse(errorMessage))
		w.Header().Add("Content-Type", "application/json")
		w.Write(resp)
		return
	}
	if !endTime.After(startTime) {
		resp, _ := json.Marshal(slack.NewErrorsViewSubmissionResponse(map[string]string{
			"end_time": "終了時刻が開始時刻よりも早いです。",
		}))
		w.Header().Add("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	submit.spreadsheetService.Add(userID, dateString, startTimeString, endTimeString)

	if _, err := submit.client.PostEphemeral(
		os.Getenv("ATTENDANCE_CHANNEL_ID"),
		userID,
		slack.MsgOptionText(message, false),
	); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
