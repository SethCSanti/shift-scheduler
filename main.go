package main

import (
    "log"
    "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from Shawn Mendix"))
}
    
func scheduleView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Schedule view"))
}

func approvalView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Approval view"))
}

func submitView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Submit view"))
}

func decideView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Decide view"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", home) // Returns full page, home page
    mux.HandleFunc("GET /schedule", scheduleView) // Returns full page, schedule view
    mux.HandleFunc("GET /approval", approvalView) // Returns full page, approval view
    mux.HandleFunc("POST /schedule/submit", submitView) // Returns HTML fragment (HTMX), submission view
    mux.HandleFunc("POST /schedule/decide", decideView) // Returns HTML fragment (HTMX), decision view

    log.Print("starting server on https://localhost:4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}