package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"wallawire/logging"
)

const (
	hContentLength = "Content-Length"
	hContentType   = "Content-Type"
	mimeTypeJson   = "application/json"
)

func sendJsonMessage(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	logger := logging.New(ctx, "sendJsonMessage")
	if len(message) == 0 {
		message = http.StatusText(statusCode)
	}
	errmsg := struct {
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message,omitempty"`
	}{
		StatusCode: statusCode,
		Message:    message,
	}
	msg, errMsg := json.Marshal(&errmsg)
	if errMsg != nil {
		logger.Error().Err(errMsg).Msg("Cannot marshal json error message")
		msg = []byte("{}")
	}
	w.Header().Set(hContentType, mimeTypeJson)
	w.Header().Set(hContentLength, strconv.Itoa(len(msg)))
	w.WriteHeader(statusCode)
	w.Write(msg)
}
