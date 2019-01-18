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
			c.CreateMessage(m.ChannelID, "Pong!")
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

			c.CreateMessageFile(m.ChannelID, f)
		} else if m.Content == "gocord embed" {
			embed := embeds.New()
			embed.SetColor("blue").SetAuthor("gocord", "").SetDescription("An awesome Golang library.")

			c.CreateMessageEmbed(m.ChannelID, embed)
		} else if m.Content == "gocord avatar" {
			avatar := m.Author.AvatarURL("", 2048)
			fmt.Println(avatar)
			embed := embeds.New()
			embed.SetAuthor(m.Author.Username, avatar)
			embed.Image = embeds.EmbedImage{
				URL: avatar,
			}

			c.CreateMessageEmbed(m.ChannelID, embed)
		}
	})

	c.Spawn()
	c.Wait()
}
