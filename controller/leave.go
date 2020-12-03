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

type Leave struct {
	client *slack.Client
	jst    *time.Location
}

func NewLeave(client *slack.Client) *Leave {
	return &Leave{
		client: client,
		jst:    time.FixedZone("Asia/Tokyo", 9*60*60),
	}
}

func (leave *Leave) HandleSlash(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.GetSlashCommandFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/out":
		userID := s.UserID
		userName := s.UserName

		// sheet.Edit(userID, time.Now().In(jst).Format("2006-01-02"), "", "", "leave")

		message := fmt.Sprintf("%s が退室しました", userName)
		if _, err := leave.client.PostEphemeral(
			os.Getenv("ATTENDANCE_CHANNEL_ID"),
			userID,
			slack.MsgOptionText(message, false),
		); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
