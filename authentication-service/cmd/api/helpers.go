package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitemtpy"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// Define a max limit of 1Mb to the request received in r.Body
	maxBytes := 1048576 //One megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Decode the request body into data
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// Check that only one message was received in the body.
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have a single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// Encode data
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// If headers are passed, add them to the response
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Send json response
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
