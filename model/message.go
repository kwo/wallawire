package model

type PushMessage struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Data string `json:"data"`
}
