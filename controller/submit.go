package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

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
	date := payload.View.State.Values["date"]["date"].SelectedDate
	startTimeString := payload.View.State.Values["start_time"]["startTime"].Value
	endTimeString := payload.View.State.Values["end_time"]["endTime"].Value

	message := fmt.Sprintf("%s が予定を追加しました\nDate: %s\nStart Time: %s\nEnd Time: %s", userName, date, startTimeString, endTimeString)

	timeRegex := regexp.MustCompile("([01][0-9]|2[0-3]):[0-5][0-9]")

	errorMessage := make(map[string]string)
	if !timeRegex.Match([]byte(startTimeString)) {
		errorMessage["start_time"] = "不正な入力です。"
	}
	if !timeRegex.Match([]byte(endTimeString)) {
		errorMessage["end_time"] = "不正な入力です。"
	}
	if len(errorMessage) != 0 {
		resp, _ := json.Marshal(slack.NewErrorsViewSubmissionResponse(errorMessage))
		w.Header().Add("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	submit.spreadsheetService.Add(userID, date, startTimeString, endTimeString)

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
