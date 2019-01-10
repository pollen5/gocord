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
	ReadyEvent = "Ready"
)

const (
	OnlinePresence       = "online"
	IdlePresence         = "idle"
	DoNotDisturbPresence = "dnd"
	InvisiblePresence    = "invisible"
)

const (
	// APIVersion is the current usable discord API version
	APIVersion = 6
	// RestURL is the base URL for all rest api requests
	RestURL     = "https://discordapp.com/api/v7"
	gatewayPath = "/gateway/bot"
)
