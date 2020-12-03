package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func VerifyingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SIGNING_SECRET"))
		if err != nil {
			fmt.Println(err.Error())
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if _, err := verifier.Write(body); err != nil {
			fmt.Println(err.Error())
		}

		if err := verifier.Ensure(); err != nil {
			fmt.Println(err.Error())
		}

		next.ServeHTTP(w, r)
	})
}
