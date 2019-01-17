package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Soumil07/gocord"
	"github.com/Soumil07/gocord/embeds"
	"github.com/Soumil07/gocord/rest"
)

func main() {
	tokenPtr := flag.String("t", "", "")
	flag.Parse()

	if *tokenPtr == "" {
		panic("No token provided.")
	}

	c := gocord.NewCluster(*tokenPtr, gocord.ClusterOptions{
		Shards: []int{0},
		Debug:  true,
	})

	c.Subscribe("ready", func(s *gocord.Shard) {
		fmt.Println("Ready to roll!")
	})
	c.Subscribe("message", func(s *gocord.Shard, m *gocord.Message) {
		if m.Content == "gocord ping" {
			s.CreateMessage(m.ChannelID, "Pong!")
		} else if m.Content == "gocord file" {
			file, err := os.Open("examples/gopher.jpg")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			f := rest.File{
				Name:        "gopher.png",
				Reader:      file,
				ContentType: "image/png",
			}

			s.CreateMessageFile(m.ChannelID, "", f)
		} else if m.Content == "gocord embed" {
			embed := embeds.New()
			embed.SetColor("blue").SetAuthor("gocord", "").SetDescription("An awesome Golang library.")

			s.CreateMessageEmbed(m.ChannelID, embed, "")
		}
	})

	c.Spawn()
	c.Wait()
}
