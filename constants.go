package gocord

// OPCodes, as documented at https://discordapp.com/developers/docs/topics/opcodes-and-status-codes
const (
	OPCodeDispatch = iota
	OPCodeHeartbeat
	OPCodeIdentify
	OPCodeStatusUpdate
	_
	OPCodeVoiceStateUpdate
	OPCodeResume
	OPCodeReconnect
	OPCodeRequestGuildMembers
	OPCodeInvalidSession
	OPCodeHello
	OPCodeHeartbeatAck
)

const (
	// ReadyEvent is the first event sent by Discord, including the session ID, guilds and other information
	ReadyEvent = "READY"
	// GuildCreateEvent is dispatched to lazy load a guild, or when a new guild is added
	GuildCreateEvent = "GUILD_CREATE"
)

const (
	OnlinePresence       = "online"
	IdlePresence         = "idle"
	DoNotDisturbPresence = "dnd"
	InvisiblePresence    = "invisible"
)

// Game activity types
const (
	// Playing gocord
	ActivityTypePlaying = iota
	// Streaming gocord
	ActivityTypeStreaming
	// Listening to gocord
	ActivityTypeListening
	// Watching gocord
	ActivityTypeWatching
)

const (
	// APIVersion is the current usable discord API version
	APIVersion = 6
	// RestURL is the base URL for all rest api requests
	RestURL     = "https://discordapp.com/api/v7"
	gatewayPath = "/gateway/bot"
)
