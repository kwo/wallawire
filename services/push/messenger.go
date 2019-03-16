package push

import (
	"sync"

	"github.com/rs/zerolog"

	"wallawire/logging"
	"wallawire/model"
)

type MessageChannel chan model.PushMessage
type SessionMap map[string]MessageChannel
type UserMap map[string]SessionMap

type PushMessenger struct {
	clientsLock          sync.RWMutex
	clients              UserMap
	onConnectTriggers    []func(userID, sessionID string)
	onDisconnectTriggers []func(userID, sessionID string)
	logger               *zerolog.Logger
}

func New() *PushMessenger {
	return &PushMessenger{
		clients: make(UserMap),
		logger:  logging.New(nil, "push"),
	}
}

func (z *PushMessenger) AddOnClientConnectTrigger(fn func(userID, sessionID string)) {
	z.clientsLock.Lock()
	defer z.clientsLock.Unlock()
	z.onConnectTriggers = append(z.onConnectTriggers, fn)
}

func (z *PushMessenger) AddOnClientDisconnectTrigger(fn func(userID, sessionID string)) {
	z.clientsLock.Lock()
	defer z.clientsLock.Unlock()
	z.onDisconnectTriggers = append(z.onDisconnectTriggers, fn)
}

func (z *PushMessenger) ConnectClient(userID, sessionID string, messageChannel chan model.PushMessage) {
	z.clientsLock.Lock()
	defer z.clientsLock.Unlock()

	sessionMap := z.clients[userID]
	if sessionMap == nil {
		sessionMap = make(SessionMap)
		z.clients[userID] = sessionMap
	}

	sessionMap[sessionID] = messageChannel

	for _, tr := range z.onConnectTriggers {
		go tr(userID, sessionID)
	}

	z.logger.Debug().Str("UserID", userID).Str("SessionID", sessionID).Msg("client connected")

}

func (z *PushMessenger) DisconnectClient(userID, sessionID string) {
	z.clientsLock.Lock()
	defer z.clientsLock.Unlock()

	sessionMap := z.clients[userID]
	if sessionMap != nil {
		if messageChan, ok := sessionMap[sessionID]; ok {
			close(messageChan)
		}
		delete(sessionMap, sessionID)
	}

	if len(sessionMap) == 0 {
		delete(z.clients, userID)
	}

	for _, tr := range z.onDisconnectTriggers {
		go tr(userID, sessionID)
	}

	z.logger.Debug().Str("UserID", userID).Str("SessionID", sessionID).Msg("client disconnected")

}

// SendMessage will send a message to all connected users if userID and sessionID are empty.
// It will send to all sessions of a specific user if sessionID is empty
// and to a specific user session if all three arguments are given.
func (z *PushMessenger) SendMessage(msg model.PushMessage, userID, sessionID string) int {

	z.clientsLock.RLock()
	defer z.clientsLock.RUnlock()

	counter := 0

	if userID == "" {
		// all
		for _, sessionMap := range z.clients {
			for _, messageChannel := range sessionMap {
				messageChannel <- msg
				counter += 1
			}
		}
		// z.logger.Debug().Interface("message", msg).Int("count", counter).Msg("sent message")
	} else if sessionID == "" {
		// all sessions for user
		if sessionMap, ok := z.clients[userID]; ok {
			for _, messageChannel := range sessionMap {
				messageChannel <- msg
				counter += 1
			}
		}
		// z.logger.Debug().Interface("message", msg).Str("UserID", userID).Int("count", counter).Msg("sent message")
	} else {
		// single session
		if sessionMap, ok := z.clients[userID]; ok {
			if messageChannel, ok2 := sessionMap[sessionID]; ok2 {
				messageChannel <- msg
				counter += 1
			}
		}
		// z.logger.Debug().Interface("message", msg).Str("UserID", userID).Str("sessionID", sessionID).Int("count", counter).Msg("sent message")
	}

	return counter

}
