package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/haydenwoodhead/gateway"
	"github.com/justinas/alice"
)

var wordTemplate = template.Must(template.ParseFiles("echoword.html"))

type IPResponse struct {
	Success bool   `json:"success"`
	IP      string `json:"ip"`
}

func main() {
	lambda, err := strconv.ParseBool(os.Getenv("LAMBDA"))

	if err != nil {
		log.Fatalf("Failed to get if lambda env var: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root"))
	}))

	r.Handle("/ip", alice.New(JSONContentType).ThenFunc(EchoIP)).Methods(http.MethodGet)
	r.HandleFunc("/echo/{word}", EchoWord).Methods(http.MethodGet)

	if lambda {
		gateway.ListenAndServe("", r)
	} else {
		log.Fatal(http.ListenAndServe(":8080", r))
	}
}

func EchoIP(w http.ResponseWriter, r *http.Request) {
	resp := IPResponse{
		Success: true,
		IP:      r.RemoteAddr,
	}

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("An error occurred"))
	}

	_, err = w.Write(jsonResp)

	if err != nil {
		log.Printf("EchoIP: failed to write response: %v", err)
	}
}

func EchoWord(w http.ResponseWriter, r *http.Request) {
	m := mux.Vars(r)
	word := m["word"]

	data := struct {
		Word string
	}{
		Word: word,
	}

	err := wordTemplate.Execute(w, data)

	if err != nil {
		log.Printf("EchoWord: failed to write template: %v", err)
	}
}

func JSONContentType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
