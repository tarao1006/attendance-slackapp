package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

type Submit struct {
	client *slack.Client
}

func NewSubmit(client *slack.Client) *Submit {
	return &Submit{client: client}
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

	// sheet.Edit(userID, date, startTime, endTime, "add")

	if _, err := submit.client.PostEphemeral(
		os.Getenv("TEST_CHANNEL_ID"),
		userID,
		slack.MsgOptionText(message, false),
	); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
