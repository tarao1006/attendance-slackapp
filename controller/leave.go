package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/sheet"
)

type Leave struct {
	client *slack.Client
}

func NewLeave(client *slack.Client) *Leave {
	return &Leave{client: client}
}

func (leave *Leave) HandleSlash(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/out":
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		userID := s.UserID
		userName := s.UserName
		sheet.Edit(userID, time.Now().In(jst).Format("2006-01-02"), "", "", "leave")
		message := fmt.Sprintf("%s が退室しました", userName)
		if _, _, err := leave.client.PostMessage(
			os.Getenv("TEST_CHANNEL_ID"),
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
