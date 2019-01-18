package gocord

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Soumil07/gocord/embeds"
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

type CreateMessage struct {
	ChannelID string
	Content   string
	Embed     *embeds.Embed
	Files     []rest.File
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
func (s *Shard) CreateMessage(channelID string, message string) (*Message, error) {
	return s.CreateMessageComplex(CreateMessage{
		ChannelID: channelID,
		Content:   message,
	})
}

func (s *Shard) CreateMessageFile(channelID string, files ...rest.File) (*Message, error) {
	return s.CreateMessageComplex(CreateMessage{
		ChannelID: channelID,
		Files:     files,
	})
}

func (s *Shard) CreateMessageEmbed(channelID string, embed *embeds.Embed) (*Message, error) {
	return s.CreateMessageComplex(CreateMessage{
		ChannelID: channelID,
		Embed:     embed,
	})
}

func (s *Shard) CreateMessageComplex(c CreateMessage) (m *Message, err error) {
	endpoint := rest.ChannelMessages(c.ChannelID)

	body, err := json.Marshal(&struct {
		Content string        `json:"content"`
		Embed   *embeds.Embed `json:"embed"`
	}{c.Content, c.Embed})
	if err != nil {
		return
	}

	err = s.Rest.Do(http.MethodPost, endpoint, body, &m, c.Files...)
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

func (s *Shard) CreateReaction(channelID, messageID, emoji string) (err error) {
	endpoint := rest.ChannelMessageReactions("@me", channelID, messageID, emoji)
	err = s.Rest.Do(http.MethodPut, endpoint, nil, nil)

	return
}

func (s *Shard) RemoveReaction(userID, channelID, messageID, emoji string) (err error) {
	endpoint := rest.ChannelMessageReactions(userID, channelID, messageID, emoji)
	err = s.Rest.Do(http.MethodDelete, endpoint, nil, nil)

	return
}

func (s *Shard) RemoveOwnReaction(channelID, messageID, emoji string) error {
	return s.RemoveReaction("@me", channelID, messageID, emoji)
}

func (s *Shard) RemoveAllReactions(channelID, messageID string) (err error) {
	endpoint := rest.ChannelMessageReactionsAll(channelID, messageID)
	err = s.Rest.Do(http.MethodDelete, endpoint, nil, nil)

	return
}

func (s *Shard) DeleteMessage(channelID, messageID string) (err error) {
	endpoint := rest.ChannelMessage(messageID, channelID)
	err = s.Rest.Do(http.MethodDelete, endpoint, nil, nil)

	return
}

func (s *Shard) BulkDeleteMessages(channelID string, amount int) (err error) {
	if amount < 2 || amount > 100 {
		return errors.New("amount must be between 2 and 100")
	}
	endpoint := rest.ChannelBulkDelete(channelID)

	body, err := json.Marshal(&struct {
		Messages int `json:"messages"`
	}{amount})
	if err != nil {
		return
	}

	err = s.Rest.Do(http.MethodPost, endpoint, body, nil)
	return
}
