package main

import (
    "log"
    "net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Shawn Mendix" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello from Shawn Mendix"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", home) // Restrict this route to exact matches on / only.
    mux.HandleFunc("GET /example", home) // This will match /example and /example/ but not /example/something.
    mux.HandleFunc("POST /example", home) // This will only match POST requests to /example.

    log.Print("starting server on https://localhost:4000")

    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}