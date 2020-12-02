package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/sheet"
)

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

func verifySigningSecret(r *http.Request) error {
	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SIGNING_SECRET"))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if _, err := verifier.Write(body); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := verifier.Ensure(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func handleSlash(w http.ResponseWriter, r *http.Request) {
	if err := verifySigningSecret(r); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(r.FormValue("payload")), &payload); err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	var userID string = payload.User.ID
	var userName string = payload.User.Name
	var date string = payload.View.State.Values["date"]["date"].SelectedDate
	var startTime string = payload.View.State.Values["start_time"]["startTime"].Value
	var endTime string = payload.View.State.Values["end_time"]["endTime"].Value

	message := fmt.Sprintf("%s が予定を追加しました\nDate: %s\nStart Time: %s\nEnd Time: %s", userName, date, startTime, endTime)

	sheet.Edit(userID, date, startTime, endTime, "add")

	api := slack.New(os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN"))
	if _, _, err := api.PostMessage(
		os.Getenv("TEST_CHANNEL_ID"),
		slack.MsgOptionText(message, false),
	); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", hello)
	http.HandleFunc("/slash", handleSlash)
	http.HandleFunc("/submit", handleSubmit)
	http.ListenAndServe(":"+port, nil)
}
