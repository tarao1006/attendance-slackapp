package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

	var userID string = payload.User.ID
	var userName string = payload.User.Name
	var date string = payload.View.State.Values["date"]["date"].SelectedDate
	var startTime string = payload.View.State.Values["start_time"]["startTime"].Value
	var endTime string = payload.View.State.Values["end_time"]["endTime"].Value

	message := fmt.Sprintf("%s が予定を追加しました\nDate: %s\nStart Time: %s\nEnd Time: %s", userName, date, startTime, endTime)

	submit.spreadsheetService.Add(userID, date, startTime, endTime)

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
