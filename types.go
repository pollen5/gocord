package gocord

import "encoding/json"

type gatewayPayload struct {
	URL    string `json:"url"`
	Shards int    `json:"shards"`
}

type receivePayload struct {
	OP  int             `json:"op"`
	D   json.RawMessage `json:"d"`
	Seq int             `json:"s,omitempty"`
	T   string          `json:"t,omitempty"`
}

type sendPayload struct {
	OP int         `json:"op"`
	D  interface{} `json:"d"`
}

type helloPayload struct {
	HeartbeatInterval int64    `json:"heartbeat_interval"`
	Trace             []string `json:"trace"`
}

type identifyPayload struct {
	Token      string             `json:"token"`
	Properties identifyProperties `json:"properties"`
	Shard      [2]int             `json:"shard"`
	Presence   Presence           `json:"presence"`
}

// Presence represents a Discord presence object
type Presence struct {
	Game   Game   `json:"game"`
	Status string `json:"status"`
}

// Game is the game status in the presence object
type Game struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type identifyProperties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type dispatch struct {
	D interface{} `json:"d"`
	T string      `json:"t"`
}

type readyDispatch struct {
	Version   string        `json:"v"`
	User      interface{}   `json:"user"`   // TODO: user type
	Guilds    []interface{} `json:"guilds"` // TODO: guild type
	SessionID string        `json:"session_id"`
}

// DispatchEvent is an event dispatched by the API
type DispatchEvent struct {
	Shard int // the ID of the shard receiving this event
	Data  interface{}
}
