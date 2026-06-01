package main

import (
    // native Go packages
    "net/http"

    // internal packages

    // 3rd party packages
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

func registerView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Register view"))
}

func registerSubmitView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Register submit view"))
}

func loginView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Login view"))
}

func loginSubmitView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Login submit view"))
}

func logoutView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Logout view"))
}

func profileView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Profile view"))
}

func pictureView(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Profile picture update view"))
}