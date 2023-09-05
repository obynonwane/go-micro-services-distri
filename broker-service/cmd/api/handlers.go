package main

import (
	"net/http"
)

// response payload type

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	//response payload
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker again",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}
