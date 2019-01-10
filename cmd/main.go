package main

// this will be updated when the lib is more complete
import (
	"flag"

	"github.com/Soumil07/gocord"
)

func main() {
	tokenPtr := flag.String("token", "", "The Discord token")
	flag.Parse()
	if *tokenPtr == "" {
		panic("Invalid token provided.")
	}

	var cluster = gocord.NewCluster(*tokenPtr, gocord.ClusterOptions{
		Shards: []int{0},
		Presence: gocord.Presence{
			Game: gocord.Game{
				Name: "gocord",
				Type: 0,
			},
			Status: gocord.DoNotDisturbPresence,
		},
	})
	cluster.Spawn()
}
