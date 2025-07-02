package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func writeJSONResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{Message: message}
	json.NewEncoder(w).Encode(response)
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "OK", http.StatusOK)
}

func error403Handler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "Forbidden", http.StatusForbidden)
}

func error500Handler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "Internal Server Error", http.StatusInternalServerError)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, "Redirecting to /redirect-two", http.StatusSeeOther)
}

func redirectTwoHandler(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	message := fmt.Sprintf("Redirected from %s", referer)
	writeJSONResponse(w, message, http.StatusOK)
}

func main() {
	http.HandleFunc("/", okHandler)
	http.HandleFunc("/error403", error403Handler)
	http.HandleFunc("/error500", error500Handler)
	http.HandleFunc("/redirect", redirectHandler)
	http.HandleFunc("/redirecttwo", redirectTwoHandler)

	log.Println("Server starting on :9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}