package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/slack-go/slack"
	"github.com/tarao1006/attendance-slackapp/controller"
	"github.com/tarao1006/attendance-slackapp/middleware"
	"github.com/tarao1006/attendance-slackapp/sheet"
)

type Server struct {
	client             *slack.Client
	router             http.Handler
	spreadsheetService *sheet.SpreadsheetService
}

func NewServer() *Server {
	return &Server{
		client:             slack.New(os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN")),
		spreadsheetService: sheet.NewSpreadsheetService(),
	}
}

func (s *Server) Init() {
	s.router = s.Route()
}

func (s *Server) Run() error {
	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, s.router)
}

func (s *Server) Route() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World!")
	})

	submitController := controller.NewSubmit(s.client)
	slackRouter := router.PathPrefix("/").Subrouter()
	slackRouter.Use(middleware.VerifyingMiddleware)
	slackRouter.HandleFunc("/submit", submitController.HandleSubmit)

	addController := controller.NewAdd(s.client)
	attendanceController := controller.NewAttendance(s.client)
	commandRouter := router.PathPrefix("/").Subrouter()
	commandRouter.Use(middleware.VerifyingMiddleware)
	commandRouter.Use(middleware.CommandMiddleware)
	commandRouter.HandleFunc("/add", addController.HandleSlash)
	commandRouter.HandleFunc("/enter", attendanceController.HandleSlash)
	commandRouter.HandleFunc("/leave", attendanceController.HandleSlash)

	return router
}
