package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)

func main() {
	http.HandleFunc("/send", SendChat)
	http.HandleFunc("/read", ReadCHat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// SendChat is endpoint to send chat
func SendChat(w http.ResponseWriter, r *http.Request) {
	untrustedPayload := r.FormValue("chat")

	trustedPayload := sanitize(untrustedPayload)
	writeToStorage(trustedPayload)

	response := fmt.Sprintf("'%s' Sent!", trustedPayload)
	fmt.Fprint(w, response)
}

// ReadCHat is endpoint to read chat
func ReadCHat(w http.ResponseWriter, r *http.Request) {
	untrustedResponse := readFromStorage()
	trustedResponse := sanitize(untrustedResponse)

	response := fmt.Sprintf("Receiving: '%s'", trustedResponse)
	fmt.Fprint(w, response)
}

func sanitize(untrustedPayload string) (trustedPayload string) {
	s := bluemonday.UGCPolicy()
	trustedPayload = s.Sanitize(untrustedPayload)
	return
}

// writeToStorage is abstraction function so write to storage DB
func writeToStorage(data string) (output bool) {
	output = true
	return
}

// readFromStorage is abstraction function to read from storage DB
func readFromStorage() (output string) {
	data := []string{
		"Good Morning",
		"Ohayou!!",
		"<a onmouseover=\"alert('XSS2')\">XSS<a>",
		"Hello <STYLE>.XSS{background-image:url(\"javascript:alert('XSS')\");}</STYLE><A CLASS=XSS></A>World",
	}

	output = data[rand.Intn(len(data))]
	return
}
