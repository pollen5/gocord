package gocord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	eventemitter "github.com/euskadi31/go-eventemitter"
)

// Cluster of Shards connecting to the gateway
type Cluster struct {
	*eventemitter.Emitter
	Token       string
	Shards      map[int]*Shard
	TotalShards int
	GatewayURL  string
	Options     ClusterOptions

	handlers sync.Map // event handlers
}

// ClusterOptions are the options used in the cluster
type ClusterOptions struct {
	Shards      []int // an array of shard IDs
	TotalShards int   // the total shards to spawn
	Presence    Presence
	Debug       bool // set to true during debug mode ONLY, this will log a lot of (useful) stuff such as reconnects and headers
}

func (c *Cluster) fetchRecommendedShards() int {
	req, _ := http.NewRequest(http.MethodGet, RestURL+gatewayPath, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var decoded gatewayPayload
	decoder.Decode(&decoded)
	c.GatewayURL = decoded.URL

	return decoded.Shards
}

// NewCluster returns a cluster instance
func NewCluster(token string, opts ClusterOptions) *Cluster {
	cluster := &Cluster{
		Emitter: eventemitter.New(),
		Token:   token,
	}
	cluster.Options = opts
	recShards := cluster.fetchRecommendedShards()

	cluster.Shards = make(map[int]*Shard)
	if len(opts.Shards) == 0 {
		if opts.TotalShards == 0 {
			totalShards := recShards
			cluster.TotalShards = totalShards
			for i := 0; i < recShards; i++ {
				opts.Shards = append(opts.Shards, i)
			}
		}
	} else {
		if opts.TotalShards == 0 {
			cluster.TotalShards = len(opts.Shards)
		}
	}

	for _, id := range opts.Shards {
		shard := NewShard(id, cluster)
		cluster.Shards[id] = shard
	}

	return cluster
}

// Spawn starts all shards and returns a slice of errors returned from every shard
func (c *Cluster) Spawn() []error {
	var wg sync.WaitGroup
	var out []error

	for _, shard := range c.Shards {
		wg.Add(1)
		err := shard.Connect()
		if err != nil {
			out = append(out, err)
		}
	}

	wg.Wait()
	return out
}

/* USEFUL SHARD-CLUSTER WRAPPERS */
func (c *Cluster) Guilds() (n int) {
	for _, shard := range c.Shards {
		n += shard.GuildCache.Size()
	}

	return
}
