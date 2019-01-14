package main

// this will be updated when the lib is more complete
import (
	"bytes"
	"flag"
	"fmt"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/Soumil07/gocord"
	humanize "github.com/dustin/go-humanize"
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
		Debug: true,
	})

	cluster.Subscribe("ready", func(s *gocord.Shard) {
		fmt.Println("Ready to roll!")
	})

	cluster.Subscribe("message", func(s *gocord.Shard, m gocord.Message) {
		isCommand := strings.HasPrefix(m.Content, "gocord")
		if !isCommand {
			return
		}
		raw := strings.SplitAfter(m.Content, "gocord")
		if len(raw) < 2 {
			return
		}
		trimmed := strings.Trim(raw[1], " ")
		rawArgs := strings.Split(trimmed, " ")
		var command string
		// var args []string

		switch len(rawArgs) {
		case 1:
			command = rawArgs[0]
		default:
			command = rawArgs[0]
			// args = rawArgs[1:]
		}

		switch command {
		case "ping":
			s.CreateMessage(m.ChannelID, "Pong!")

		case "stats":
			stats := &runtime.MemStats{}
			runtime.ReadMemStats(stats)

			buf := &bytes.Buffer{}
			tab := &tabwriter.Writer{}
			tab.Init(buf, 0, 0, 4, ' ', 0)

			fmt.Fprintf(tab, "```\n")
			fmt.Fprintf(tab, "Gocord: \t%s\n", gocord.VERSION)
			fmt.Fprintf(tab, "Go: \t%s\n", runtime.Version())
			fmt.Fprintf(tab, "Memory used: \t%s / %s (%s garbage collected)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
			fmt.Fprintf(tab, "Guilds: \t%d", cluster.Guilds())
			fmt.Fprintf(tab, "```")

			tab.Flush()
			out := buf.String()

			s.CreateMessage(m.ChannelID, out)
		}
	})

	cluster.Spawn()
	cluster.Wait()
}
