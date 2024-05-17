package handler

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func writeResponseJson(w http.ResponseWriter, body interface{}) {
	resp, err := jsoniter.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
