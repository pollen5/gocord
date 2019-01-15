package gocord

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/Soumil07/gocord/cache"
	"github.com/Soumil07/gocord/rest"
	"github.com/gorilla/websocket"
)

// Shard represents a Shard connecting to the gateway. All underlying WS connections
// are done through Shards, with events being forwarded to the main Cluster
type Shard struct {
	sync.RWMutex
	Cluster           *Cluster
	ws                *websocket.Conn
	heartbeatTicker   *time.Ticker
	lastHeartbeatSent int64
	heartbeatAcked    bool // whether the heartbeat has been acknowledged

	Latency int64 // heartbeat ack latency

	ID         int
	Token      string
	Seq        int
	SessionID  string
	GuildCache *cache.Cache // a mutable LRU cache with capacity set to 0
	Rest       *rest.RestManager
}

// NewShard returns a new shard instance
func NewShard(ID int, cluster *Cluster) *Shard {
	shard := &Shard{
		Cluster:        cluster,
		heartbeatAcked: true,

		ID:         ID,
		Token:      cluster.Token,
		GuildCache: cache.NewCache(0),
		Rest:       rest.NewRestManager(cluster.Token),
	}

	return shard
}

// Connect establishes a connection with the Discord API
func (s *Shard) Connect() (err error) {
	// TODO: forward this to the rest API when done
	s.Cluster.GatewayURL = fmt.Sprintf("%s?v=%d&encoding=json", s.Cluster.GatewayURL, APIVersion)
	s.ws, _, err = websocket.DefaultDialer.Dial(s.Cluster.GatewayURL, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to gateway: %s", err.Error()))
	}
	s.ws.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	msgChan, errChan := s.listen()
	for {
		select {
		case m := <-msgChan:
			s.onMessage(m)
		case err := <-errChan:
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			}
		}
	}
}

// wrapper around sending to WS
func (s *Shard) send(op int, d interface{}) error {
	s.Lock()
	defer s.Unlock()

	return s.ws.WriteJSON(&sendPayload{
		OP: op,
		D:  d,
	})
}

func (s *Shard) onMessage(packet *receivePayload) error {
	// update the last sequence received
	if packet.Seq != 0 {
		s.Seq = packet.Seq
	}

	switch packet.OP {
	case OPCodeHello:
		var pk helloPayload
		err := json.Unmarshal(packet.D, &pk)
		if err != nil {
			return err
		}

		s.debugf("heartbeat interval: %d", pk.HeartbeatInterval)
		go s.startHeartbeat(time.Duration(pk.HeartbeatInterval) * time.Millisecond)

		s.debug("sending identify payload")
		return s.identify()

	case OPCodeDispatch:
		switch packet.T {
		case ReadyEvent:
			var pk readyDispatch
			err := json.Unmarshal(packet.D, &pk)
			if err != nil {
				panic(err)
			}

			var unavailableGuilds int
			for _, guild := range pk.Guilds {
				if guild.Unavailable {
					unavailableGuilds++
				}
				s.GuildCache.Add(guild.ID, guild)
			}
			s.debugf("%d guilds loaded, %d unavailable", s.GuildCache.Size(), unavailableGuilds)

			s.Cluster.Dispatch("ready", s)

		// GUILD_CREATE is sometimes fired immediately after ready to load all lazy loaded guilds
		case GuildCreateEvent:
			var guild Guild
			err := json.Unmarshal(packet.D, &guild)
			if err != nil {
				panic(err)
			}

			// lazy loading unavailable guilds, don't dispatch GUILD_CREATE to the cluster
			if s.GuildCache.Has(guild.ID) {
				s.debugf("lazy loaded the guild %s", guild.Name)
				s.GuildCache.Update(guild.ID, guild)
			} else {
				s.GuildCache.Add(guild.ID, guild)
				s.Cluster.Dispatch("guildCreate", guild)
			}

		case MessageEvent:
			var m Message
			err := json.Unmarshal(packet.D, &m)
			if err != nil {
				panic(err)
			}

			s.Cluster.Dispatch("message", s, *m)
		}

	case OPCodeHeartbeatAck:
		s.Latency = (time.Now().UnixNano() - s.lastHeartbeatSent) / 1000000 // nanoseconds to milliseconds
		s.debugf("heartbeat acknowledged. latency: %d ms", s.Latency)
		s.heartbeatAcked = true
	}

	return nil
}

// UpdatePresence updates the shard's presence. NOTE: check UpdateGame if you just want to set the game
func (s *Shard) UpdatePresence(presence Presence) error {
	return s.send(OPCodeStatusUpdate, presence)
}

// UpdateGame updates the shard's presence to the given game
func (s *Shard) UpdateGame(game string) error {
	presence := Presence{
		Status: OnlinePresence,
		Game: Game{
			Name: game,
			Type: ActivityTypePlaying,
		},
	}

	return s.UpdatePresence(presence)
}

// UpdateGamef is UpdateGame with a format string
func (s *Shard) UpdateGamef(format string, a ...interface{}) error {
	return s.UpdateGame(fmt.Sprintf(format, a...))
}

func (s *Shard) identify() error {
	return s.send(OPCodeIdentify, identifyPayload{
		Token: s.Token,
		Properties: identifyProperties{
			OS:      runtime.GOOS,
			Browser: "gocord",
			Device:  "gocord",
		},
		Shard:          [2]int{s.ID, s.Cluster.TotalShards},
		Presence:       s.Cluster.Options.Presence,
		LargeThreshold: 250,
	})
}

func (s *Shard) Resume() error {
	s.debug("resuming connection to WS")
	return s.send(OPCodeResume, &resumeDispatch{
		Token:     s.Token,
		Sequence:  s.Seq,
		SessionID: s.SessionID,
	})
}

func (s *Shard) listen() (<-chan *receivePayload, <-chan error) {
	msgChan, errChan := make(chan *receivePayload), make(chan error)

	go func() {
		for {
			payload := &receivePayload{}
			err := s.ws.ReadJSON(payload)
			if err != nil {
				errChan <- err
				break
			}

			msgChan <- payload
		}
	}()

	return msgChan, errChan
}

func (s *Shard) startHeartbeat(duration time.Duration) {
	s.heartbeatTicker = time.NewTicker(duration)

	// call it once because thats how *time.Ticker works. not start ->work -> wait -> work, just
	// start -> wait -> work
	s.heartbeat()
	for range s.heartbeatTicker.C {
		s.heartbeat()
	}
}

// used to log stuff for debugging
func (s *Shard) debug(txt string) {
	if s.Cluster.Options.Debug {
		fmt.Printf("[%d] DEBUG: %s\n", s.ID, txt)
	}
}

func (s *Shard) debugf(format string, a ...interface{}) {
	s.debug(fmt.Sprintf(format, a...))
}

func (s *Shard) heartbeat() error {
	if !s.heartbeatAcked {
		// heartbeat hasn't been acknowledged. close the connection and attempt to resume
		err := s.Close()
		if err != nil {
			return err
		}

		s.debug("heartbeat not acknowledged, attempting a reconnect")
		return s.Resume()
	}
	err := s.send(OPCodeHeartbeat, s.Seq)
	s.debug("heartbeat sent")
	// sent heartbeat hasn't been acknowledged yet
	s.heartbeatAcked = false
	// this is to track ack latency
	s.lastHeartbeatSent = time.Now().UnixNano()
	return err
}

// Close gracefully closes the connection to Discord
func (s *Shard) Close() error {
	err := s.ws.Close()
	if err != nil {
		return err
	}
	s.heartbeatTicker.Stop()

	return nil
}
