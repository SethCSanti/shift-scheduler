package main

import (
	// native Go packages
	"log"
	"net/http"
	// internal packages
	// 3rd party packages
)

var app AppState

func main() {
	app.Users = make(map[string]*User)
	app.Schedule = make(map[string]*Schedule)
	loadData()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /{$}", home)                     // Returns full page, home page
	mux.HandleFunc("GET /schedule", scheduleView)        // Returns full page, schedule view
	mux.HandleFunc("GET /approval", approvalView)        // Returns full page, approval view
	mux.HandleFunc("POST /schedule/submit", submitView)  // Returns HTML fragment (HTMX), submission view
	mux.HandleFunc("POST /schedule/decide", decideView)  // Returns HTML fragment (HTMX), decision view
	mux.HandleFunc("GET /register", registerView)        // Returns HTML fragment (HTMX), registration view
	mux.HandleFunc("POST /register", registerSubmitView) // Returns HTML fragment (HTMX), registration submission view
	mux.HandleFunc("GET /login", loginView)              // Returns HTML fragment (HTMX), login view
	mux.HandleFunc("POST /login", loginSubmitView)       // Returns HTML fragment (HTMX), login submission view
	mux.HandleFunc("POST /logout", logoutView)           // Returns HTML fragment (HTMX), logout view
	mux.HandleFunc("GET /profile", profileView)          // Returns full page, profile view
	mux.HandleFunc("POST /profile/picture", pictureView) // Returns HTML fragment (HTMX), profile picture update view

	log.Print("starting server on http://localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
