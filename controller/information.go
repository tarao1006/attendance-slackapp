package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/httputil"
	"github.com/tarao1006/attendance-slackapp/sheet"
)

type Information struct {
	client             *slack.Client
	spreadsheetService *sheet.SpreadsheetService
	jst                *time.Location
}

func NewInformation(client *slack.Client, spreadsheetService *sheet.SpreadsheetService) *Information {
	return &Information{
		client:             client,
		spreadsheetService: spreadsheetService,
		jst:                time.FixedZone("Asia/Tokyo", 9*60*60),
	}
}

func (information *Information) HandleSlash(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.GetSlashCommandFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/info":
		dateString := s.Text
		_, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			if _, err := information.client.PostEphemeral(
				os.Getenv("ATTENDANCE_CHANNEL_ID"),
				s.UserID,
				slack.MsgOptionText("コマンドが不正です。例) /info 2006-01-02", false),
			); err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			// information.spreadsheetService.GetInformation(dateString)
			message := "info"
			if _, err := information.client.PostEphemeral(
				os.Getenv("ATTENDANCE_CHANNEL_ID"),
				s.UserID,
				slack.MsgOptionText(message, false),
			); err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
