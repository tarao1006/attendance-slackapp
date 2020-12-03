package controller

import (
	"fmt"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/httputil"
)

type Add struct {
	client *slack.Client
}

func NewAdd(client *slack.Client) *Add {
	return &Add{client: client}
}

func (add *Add) HandleSlash(w http.ResponseWriter, r *http.Request) {
	s, err := httputil.GetSlashCommandFromContext(r.Context())
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch s.Command {
	case "/add":
		modalRequest := generateModalRequest()
		_, err = add.client.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
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
