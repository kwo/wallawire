package sse

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"wallawire/logging"
	"wallawire/model"
)

const (
	hCacheControl       = "Cache-Control"
	hConnection         = "Connection"
	hContentType        = "Content-Type"
	cacheNoCache        = "no-cache"
	connectionKeepAlive = "keep-alive"
	mimetypeEventStream = "text/event-stream"
)

type PushMessenger interface {
	ConnectClient(userID, sessionID string, client chan model.PushMessage)
	DisconnectClient(userID, sessionID string)
}

func Handler(pushMessenger PushMessenger) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		logger := logging.New(ctx, "web", "sse")
		token := model.TokenFromContext(ctx)

		// Make sure that the writer supports flushing.
		flusher, ok := w.(http.Flusher)
		if !ok {
			msg := "streaming unsupported"
			logger.Warn().Msg(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set(hContentType, mimetypeEventStream)
		w.Header().Set(hCacheControl, cacheNoCache)
		w.Header().Set(hConnection, connectionKeepAlive)

		// define the message channel
		messageChan := make(chan model.PushMessage)

		// add message channel to message pushMessenger
		pushMessenger.ConnectClient(token.ID, token.SessionID, messageChan)
		var closed int32

		// failsafe: remove connection upon exit of handler
		// probably already called on request close (below) when this is triggered
		defer func() {
			if isClosed := atomic.LoadInt32(&closed); isClosed == 0 {
				pushMessenger.DisconnectClient(token.ID, token.SessionID)
			}
		}()

		// clean close: remove connection when client closes request
		go func() {
			select {
			case <-ctx.Done():
				pushMessenger.DisconnectClient(token.ID, token.SessionID)
				atomic.AddInt32(&closed, 1)
			}
		}()

		// if _, err := fmt.Fprintf(w, "retry: %d\n\n", 3000); err != nil {
		// 	logger.Error().Err(err).Msg("error sending message to client")
		// }

		// read until channel is closed
		// pushMessenger.CloseClient will close the messageChannel to break out of the loop
		for msg := range messageChan {

			// see https://hpbn.co/server-sent-events-sse/#event-stream-protocol
			if len(msg.ID) != 0 {
				if _, err := fmt.Fprintf(w, "id: %s\n", msg.ID); err != nil {
					logger.Error().Err(err).Interface("message", msg).Msg("error sending message id client")
					continue
				}
			}
			if len(msg.Type) != 0 {
				if _, err := fmt.Fprintf(w, "event: %s\n", msg.Type); err != nil {
					logger.Error().Err(err).Interface("message", msg).Msg("error sending message type to client")
					continue
				}
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", msg.Data); err != nil {
				logger.Error().Err(err).Interface("message", msg).Msg("error sending message to client")
			}

			// Flush the data immediately instead of buffering it for later.
			flusher.Flush()
		}

	})

}
