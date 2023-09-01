package main

import (
	"encoding/json"
	"net/http"
)

// response payload type

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	//response payload
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	//construct response
	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}
