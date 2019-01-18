package rest

import "fmt"

// implements helper functions for REST API endpoints

func ChannelMessages(channelID string) string {
	return format("/channels/%s/messages", channelID)
}

func ChannelMessage(messageID, channelID string) string {
	return format("/channels/%s/messages/%s", channelID, messageID)
}

func ChannelMessageReactions(userID, channelID, messageID, emoji string) string {
	return format("%s/reactions/%s/%s", ChannelMessage(channelID, messageID), emoji, userID)
}

func ChannelMessageReactionsAll(channelID, messageID string) string {
	return format("/channels/%s/messages/%s/reactions", channelID, messageID)
}

func ChannelBulkDelete(channelID string) string {
	return format("%s/bulk-delete", ChannelMessages(channelID))
}

func User(ID string) string {
	return format("/users/%s", ID)
}

func UserGuild(ID, guildID string) string {
	return format("/users/%s/guilds/%s", ID, guildID)
}

func GuildBanMember(guildID, userID string) string {
	return format("/guilds/%s/bans/%s", guildID, userID)
}

func Invite(code string) string {
	return format("/invites/%s", code)
}

func format(text string, a ...interface{}) string {
	return fmt.Sprintf(text, a...)
}
