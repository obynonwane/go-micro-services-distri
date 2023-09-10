package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// custom json response type
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// function to read json
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxByte := 1048576 //one megabyte

	//check maximum byte sent - by reading the request body
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))

	//Create a new json decorder - initialized with the json bode
	dec := json.NewDecoder(r.Body)
	//Does the actual json decoding & parsing
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	//decode the JSON data from the request body into an empty anonymous struct
	err = dec.Decode(&struct{}{})

	//if error is any other error not equal to io.EOF
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// method to write json with a variadic parameter that can receive zero or more parameter (...http.Header)
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	//converts the passed data into json representative
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	//check if any header is supplied and set the respnse header
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
