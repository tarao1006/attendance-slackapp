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

type Attendance struct {
	client             *slack.Client
	spreadsheetService *sheet.SpreadsheetService
	jst                *time.Location
}

func NewAttendance(client *slack.Client, spreadsheetService *sheet.SpreadsheetService) *Attendance {
	return &Attendance{
		client:             client,
		spreadsheetService: spreadsheetService,
		jst:                time.FixedZone("Asia/Tokyo", 9*60*60),
	}
}

func (attendance *Attendance) HandleSlash(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.GetSlashCommandFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/in", "/out":
		userID := s.UserID
		userName := s.UserName
		var t string
		if s.Command == "/in" {
			t = "入室"
			attendance.spreadsheetService.Enter(userID)
		} else {
			t = "退室"
			attendance.spreadsheetService.Leave(userID)
		}
		message := fmt.Sprintf("%s が%sしました", userName, t)
		if _, _, err := attendance.client.PostMessage(
			os.Getenv("ATTENDANCE_CHANNEL_ID"),
			slack.MsgOptionText(message, false),
		); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
