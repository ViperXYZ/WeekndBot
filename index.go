package WeekndBot

import (
	"encoding/json"
	"log"
	"net/http"
)

type testStruct struct {
	Test string
}

func process(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t testStruct
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatalf("Could not deserialize JSON request to struct: %v", err)
	}
	log.Println(t.Test)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		log.Fatalf("Could not serialize and send repsonse struct: %v", err)
	}
}
