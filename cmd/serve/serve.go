package serve

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jdkato/prose/v2"
	"github.com/mpvl/unique"
)

func tokenize(w http.ResponseWriter, r *http.Request) {
	var tokenSlice []string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	doc, err := prose.NewDocument(string(body))
	for _, tok := range doc.Tokens() {
		tokenSlice = append(tokenSlice, tok.Text)
	}

	unique.Sort(unique.StringSlice{&tokenSlice})
	jArr, err := json.Marshal(tokenSlice)
	if err != nil {
		http.Error(w, "could not parse request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jArr)
}

func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/tokenize", tokenize).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", router))
}
