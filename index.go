package WeekndBot

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/nlopes/slack"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	signingSecret, ok := os.LookupEnv("SIGNING_SECRET")
	if !ok {
		log.Println("SIGNING_SECRET not declared")
		return
	}
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Printf("Could not return SecretsVerifier object: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("Could not parse slash command: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Printf("Could not verify secret: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/weekndsays":
		params := &slack.Msg{Text: s.Text}
		b, err := json.Marshal(params)
		if err != nil {
			log.Printf("Could not marshal json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
