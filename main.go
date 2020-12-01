package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func generateModalRequest() slack.ModalViewRequest {
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

	var modalRequest slack.ModalViewRequest
	modalRequest.Type = slack.ViewType("modal")
	modalRequest.Title = titleText
	modalRequest.Close = closeText
	modalRequest.Submit = submitText
	modalRequest.Blocks = blocks
	return modalRequest
}

func verifySigningSecret(r *http.Request) error {
	signingSecret := os.Getenv("SIGNING_SECRET")
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
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

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func handleSlash(w http.ResponseWriter, r *http.Request) {
	err := verifySigningSecret(r)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	switch s.Command {
	case "/test":
		api := slack.New(os.Getenv("OAUTH_ACCESS_TOKEN"))
		modalRequest := generateModalRequest()
		_, err = api.OpenView(s.TriggerID, modalRequest)
		if err != nil {
			fmt.Printf("Error opening view: %s", err)
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

	log.Println(payload.Type == slack.InteractionType("view_submission"))
	log.Printf("Channel ID: %s\n", payload.Channel.ID)
	log.Printf("User    ID: %s\n", payload.User.ID)
	log.Printf("Hash      : %s", payload.ViewSubmissionCallback.Hash)

	log.Println(payload.View.State.Values)

	for _, v := range payload.View.State.Values {
		for k, vv := range v {
			if k == "date" {
				log.Printf("Key : %s   Value : %v", k, vv.SelectedDate)
			} else if k == "startTime" || k == "endTime" {
				log.Printf("Key : %s   Value : %v", k, vv.Value)
			}
		}
	}

	api := slack.New(os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN"))
	_, _, err := api.PostMessage(
		os.Getenv("TEST_CHANNEL_ID"),
		slack.MsgOptionText("Succeeded", false),
		slack.MsgOptionAsUser(false),
	)
	if err != nil {
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
	http.HandleFunc("/test", handleSlash)
	http.HandleFunc("/submit", handleSubmit)
	http.ListenAndServe(":"+port, nil)
}
