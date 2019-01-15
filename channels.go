package gocord

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Soumil07/gocord/rest"
)

// contains channel related structs and methods

type MessageType int

// Message types, as documented at https://discordapp.com/developers/docs/resources/channel#message-object-message-types
const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
)

// Channel represents a generic Discord channel
type Channel struct {
}

type Message struct {
	ID              string `json:"id"`
	ChannelID       string `json:"channel_id"`
	GuildID         string `json:"guild_id,omitempty"`
	Author          User   `json:"author,omitempty"`
	Content         string `json:"content"`
	Timestamp       string `json:"timestamp"`
	EditedTimestamp string `json:"edited_timestamp"`
	TTS             bool   `json:"tts"`
	MentionEveryone bool   `json:"mention_everyone"`
}

// CreatedAt returns a time object representing when the message was created
func (m *Message) CreatedAt() (time.Time, error) {
	return time.Parse(time.RFC3339, m.Timestamp)
}

// EditedAt returns a time object representing when the message was edited
func (m *Message) EditedAt() (time.Time, error) {
	if m.EditedTimestamp == "" {
		return time.Time{}, errors.New("the provided message hasn't been edited yet")
	}
	return time.Parse(time.RFC3339, m.EditedTimestamp)
}

// CreateMessage sends a message to the specified channel
func (s *Shard) CreateMessage(channelID string, message string) (m *Message, err error) {
	endpoint := rest.ChannelMessages(channelID)

	body, err := json.Marshal(&struct {
		Content string `json:"content"`
	}{message})
	if err != nil {
		return
	}

	err = s.Rest.Do(http.MethodPost, endpoint, body, &m)
	if err != nil {
		return
	}

	return
}

func (s *Shard) EditMessage(channelID, messageID, message string) (m *Message) {
	endpoint := rest.ChannelMessage(messageID, channelID)

	body, err := json.Marshal(&struct {
		Content string `json:"content"`
	}{message})
	if err != nil {
		return
	}

	err = s.Rest.Do(http.MethodPatch, endpoint, body, &m)
	if err != nil {
		return
	}

	return
}
