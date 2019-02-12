package status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wallawire/model"
)

const (
	hContentLength = "Content-Length"
	hContentType   = "Content-Type"
	mimeTypeJson   = "application/json"
	mimeTypeText   = "text/plain; charset=utf8"
)

func Handler(status *model.Status) (http.HandlerFunc, error) {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status.PopulateNow()

		payload, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			sendMessage(w, http.StatusInternalServerError)
		}

		w.Header().Set(hContentType, mimeTypeJson)
		w.Header().Set(hContentLength, strconv.Itoa(len(payload)))
		w.Write(payload)
	}), nil

}

func sendMessage(w http.ResponseWriter, statusCode int) {
	msg := []byte(http.StatusText(statusCode) + "\n")
	w.Header().Set(hContentType, mimeTypeText)
	w.Header().Set(hContentLength, strconv.Itoa(len(msg)))
	w.WriteHeader(statusCode)
	w.Write(msg)
}
