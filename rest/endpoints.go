package rest

import "fmt"

// implements helper functions for REST API endpoints

func ChannelMessages(channelID string) string {
	return format("/channels/%s/messages", channelID)
}

func ChannelMessage(messageID, channelID string) string {
	return format("/channels/%s/messages/%s", channelID, messageID)
}

func format(text string, a ...interface{}) string {
	return fmt.Sprintf(text, a...)
}
