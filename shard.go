package gocord

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/Soumil07/gocord/cache"
	"github.com/gorilla/websocket"
)

// Shard represents a Shard connecting to the gateway. All underlying WS connections
// are done through Shards, with events being forwarded to the main Cluster
type Shard struct {
	sync.Mutex
	cluster         *Cluster
	ws              *websocket.Conn
	heartbeatTicker *time.Ticker
	heartbeatAcked  bool // whether the heartbeat has been acknowledged

	ID         int
	Token      string
	Seq        int
	SessionID  string
	GuildCache *cache.Cache // a mutable LRU cache with capacity set to 0
}

// NewShard returns a new shard instance
func NewShard(ID int, cluster *Cluster) *Shard {
	shard := &Shard{
		Mutex:          sync.Mutex{},
		cluster:        cluster,
		heartbeatAcked: true,

		ID:    ID,
		Token: cluster.Token,
	}

	return shard
}

// Connect establishes a connection with the Discord API
func (s *Shard) Connect() (err error) {
	s.Lock()
	defer s.Unlock()

	// TODO: forward this to the rest API when done
	s.cluster.GatewayURL = fmt.Sprintf("%s?v=%d&encoding=json", s.cluster.GatewayURL, APIVersion)
	s.ws, _, err = websocket.DefaultDialer.Dial(s.cluster.GatewayURL, nil)
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

		go s.startHeartbeat(time.Duration(pk.HeartbeatInterval))
		return s.identify()

	case OPCodeDispatch:
		switch packet.T {
		case ReadyEvent:
			var pk readyDispatch
			err := json.Unmarshal(packet.D, &pk)
			if err != nil {
				return err
			}
			// TODO: cache guilds
		}
	}

	return nil
}

func (s *Shard) identify() error {
	return s.send(OPCodeIdentify, identifyPayload{
		Token: s.Token,
		Properties: identifyProperties{
			OS:      runtime.GOOS,
			Browser: "gocord",
			Device:  "gocord",
		},
		Shard:    [2]int{s.ID, s.cluster.TotalShards},
		Presence: s.cluster.Options.Presence,
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

	for range s.heartbeatTicker.C {
		s.heartbeat()
	}
}

func (s *Shard) heartbeat() error {
	// TODO: disconnect if heartneatAcked = false
	err := s.send(OPCodeHeartbeat, s.Seq)
	// sent heartbeat hasn't been acknowledged yet
	s.heartbeatAcked = false
	return err
}
