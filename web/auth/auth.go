package auth

import (
	"net/http"
	"strconv"
)

const (
	CookieName = "jwt"
)

const (
	hContentLength = "Content-Length"
	hContentType   = "Content-Type"
	mimeTypeJson   = "application/json"
)

func sendMessage(w http.ResponseWriter, statusCode int) {
	sendMessageText(w, statusCode, http.StatusText(statusCode))
}

func sendMessageText(w http.ResponseWriter, statusCode int, statusText string) {
	msg := []byte(statusText + "\n")
	w.Header().Set(hContentLength, strconv.Itoa(len(msg)))
	w.WriteHeader(statusCode)
	w.Write(msg)
}
