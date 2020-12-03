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

func HandleSlash(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/add":
		api := slack.New(os.Getenv("OAUTH_ACCESS_TOKEN"))
		modalRequest := generateModalRequest()
		_, err = api.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
		}
	case "/in":
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		userID := s.UserID
		userName := s.UserName
		log.Printf("%s %s", userID, userName)
		sheet.Edit(userID, time.Now().In(jst).Format("2006-01-02"), "", "", "enter")
		message := fmt.Sprintf("%s が入室しました", userName)
		api := slack.New(os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN"))
		if _, _, err := api.PostMessage(
			os.Getenv("TEST_CHANNEL_ID"),
			slack.MsgOptionText(message, false),
		); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/out":
		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		userID := s.UserID
		userName := s.UserName
		sheet.Edit(userID, time.Now().In(jst).Format("2006-01-02"), "", "", "leave")
		message := fmt.Sprintf("%s が退室しました", userName)
		api := slack.New(os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN"))
		if _, _, err := api.PostMessage(
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

func generateModalRequest() (modalRequest slack.ModalViewRequest) {
	titleText := slack.NewTextBlockObject("plain_text", "出席管理App", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Close", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "滞在予定時刻を入力してください", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	datePickerText := slack.NewTextBlockObject("plain_text", "日にち", false, false)
	datePickerElement := slack.NewDatePickerBlockElement("date")
	datePicker := slack.NewInputBlock("date", datePickerText, datePickerElement)

	startTimeText := slack.NewTextBlockObject("plain_text", "開始時刻", false, false)
	startTimePlaceholder := slack.NewTextBlockObject("plain_text", "例) 12:00", false, false)
	startTimeElement := slack.NewPlainTextInputBlockElement(startTimePlaceholder, "startTime")
	startTime := slack.NewInputBlock("start_time", startTimeText, startTimeElement)

	endTimeText := slack.NewTextBlockObject("plain_text", "終了時刻", false, false)
	endTimePlaceholder := slack.NewTextBlockObject("plain_text", "例) 12:00", false, false)
	endTimeElement := slack.NewPlainTextInputBlockElement(endTimePlaceholder, "endTime")
	endTime := slack.NewInputBlock("end_time", endTimeText, endTimeElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			datePicker,
			startTime,
			endTime,
		},
	}

	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	modalRequest.ClearOnClose = true

	return modalRequest
}
