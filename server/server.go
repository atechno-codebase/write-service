package server

import (
	"net/http"
	"write/config"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var PORT string

func Start() {
	PORT = config.Configuration.Port
	r := mux.NewRouter()

	initRoutes(r)

	logrus.Infof("started server on port: %s", PORT)
	logrus.Fatal(http.ListenAndServe(":"+PORT, r))
}

func initRoutes(r *mux.Router) {
	r.Use(
		// middleware.Logger,
		middleware.RealIP,
		middleware.RequestID,
		otelmux.Middleware("write-service"))

	// r.Use(func(h http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Println(r.Context())
	// 		h.ServeHTTP(w, r)
	// 	})
	// })

	r.HandleFunc("/reading", WriteReading).
		Methods(http.MethodPost)

	r.HandleFunc("/setpoints/{uid}", GetSetpoints).
		Methods(http.MethodGet)
}
