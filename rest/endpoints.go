package rest

import "fmt"

// implements helper functions for REST API endpoints

func ChannelMessage(channelID string) string {
	return format("/channels/%s/messages", channelID)
}

func format(text string, a ...interface{}) string {
	return fmt.Sprintf(text, a...)
}
