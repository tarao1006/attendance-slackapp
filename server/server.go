package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/controller"
)

type Server struct {
	client *slack.Client
	router http.Handler
}

func NewServer() *Server {
	return &Server{
		client: slack.New(os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN")),
		router: Route(),
	}
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, s.router)
}

func Route() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World!")
	})
	r.HandleFunc("/slash", controller.HandleSlash)
	r.HandleFunc("/submit", controller.HandleSubmit)
	r.Use(VerifyingMiddleware)
	return r
}
