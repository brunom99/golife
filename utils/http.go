package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func HttpBadRequest(w http.ResponseWriter, msg string) {
	log.Printf("StatusBadRequest: %s", msg)
	w.WriteHeader(http.StatusBadRequest)
}

func HttpError(w http.ResponseWriter, errs ...error) {
	w.WriteHeader(http.StatusInternalServerError)
	if len(errs) > 0 {
		log.Printf("StatusInternalServerError: %s", errs[0].Error())
	}
}

func HttpAccept(w http.ResponseWriter, i interface{}) {
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(i)
	_, _ = w.Write(data)
}
