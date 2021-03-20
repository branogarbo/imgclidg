package main

import (
	"flag"
	"log"

	dg "github.com/bwmarrin/discordgo"
)

var (
	token   string
	session *dg.Session
	err     error
)

func init() {
	flag.StringVar(&token, "t", "", "Bot access token.")
	flag.Parse()
}

func main() {
	session, err = dg.New("Bot " + token)
	if err != nil {
		log.Fatal("failed to start bot:", err)
	}

	session.Open()
	defer session.Close()

	session.AddHandler(func(s *dg.Session, e *dg.MessageCreate) {
		var (
			msg string = e.Message.Content
		)

		_, err = s.ChannelMessageSend("imgcli", msg)
		if err != nil {
			log.Fatal("failed to send message:", err)
		}

	})
}
