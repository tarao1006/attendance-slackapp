package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/httputil"
)

type Enter struct {
	client *slack.Client
	jst    *time.Location
}

func NewEnter(client *slack.Client) *Enter {
	return &Enter{
		client: client,
		jst:    time.FixedZone("Asia/Tokyo", 9*60*60),
	}
}

func (enter *Enter) HandleSlash(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.GetSlashCommandFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/in":
		userID := s.UserID
		userName := s.UserName
		log.Printf("%s %s", userID, userName)

		// sheet.Edit(userID, time.Now().In(jst).Format("2006-01-02"), "", "", "enter")

		message := fmt.Sprintf("%s が入室しました", userName)
		if _, err := enter.client.PostEphemeral(
			os.Getenv("ATTENDANCE_CHANNEL_ID"),
			userID,
			slack.MsgOptionText(message, false),
		); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
