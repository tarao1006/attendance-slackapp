package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

type contextKey string

const commandKey contextKey = "command"

func CommandMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), commandKey, s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
