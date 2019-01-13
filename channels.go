package gocord

import (
	"encoding/json"
	"net/http"

	"github.com/Soumil07/gocord/rest"
)

// contains channel related structs and methods

// Channel represents a generic Discord channel
type Channel struct {
}

func (s *Shard) CreateMessage(channelID string, message string) {
	endpoint := rest.ChannelMessage(channelID)

	body, err := json.Marshal(&struct {
		Content string `json:"content"`
	}{Content: message})
	if err != nil {
		panic(err)
	}

	err = s.Rest.Do(http.MethodPost, endpoint, body, &struct{}{})
	if err != nil {
		panic(err)
	}
}
