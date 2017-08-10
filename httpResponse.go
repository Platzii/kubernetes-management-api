package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// HTTPResponse represents a HTTP response
type HTTPResponse struct {
	Code    int         `json:"code"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func (r *HTTPResponse) send(w http.ResponseWriter) error {
	json, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("build body: %s", err.Error())
	}

	w.WriteHeader(r.Code)
	_, err = w.Write(json)
	if err != nil {
		return fmt.Errorf("write body: %s", err.Error())
	}
	return nil
}

func returnError(w http.ResponseWriter, code int, error, message string) {
	r := &HTTPResponse{
		Code:    code,
		Error:   error,
		Message: message,
	}
	if err := r.send(w); err != nil {
		logrus.Errorf("Could not send response: %s", err.Error())
	}
}

func returnResult(w http.ResponseWriter, result interface{}) {
	r := &HTTPResponse{
		Code:    200,
		Error:   "",
		Message: "OK",
		Result:  result,
	}

	if err := r.send(w); err != nil {
		logrus.Errorf("Could not send response: %s", err.Error())
		return
	}
}
