package middleware

import (
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/httputil"
)

func CommandMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := httputil.SetSlashCommandToContext(r.Context(), &s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
