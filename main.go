package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	imgcli "github.com/branogarbo/imgcli/util"
	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		msgHead     string
		imgUrl      string
		pixelString string
		errMsg      string
		err         error
	)

	if m.Author.ID == s.State.User.ID || len(m.Content) < 5 {
		return
	}

	// using janky command code for now

	msgHead = m.Content[:5]

	if msgHead == "!img " {
		if m.Content == msgHead {
			err = errors.New("url not provided")
		} else {
			imgUrl = m.Content[5:]

			pixelString, err = imgcli.OutputImage(imgcli.OutputConfig{
				Src:          strings.TrimSpace(imgUrl),
				OutputMode:   "ascii",
				AsciiPattern: " .:-=+*#%@",
				IsUseWeb:     true,
				OutputWidth:  66,
			})
		}

		if err != nil {
			errMsg = fmt.Sprintf("imgcli/util: failed to output image: %v", err)

			fmt.Println(errMsg)
			pixelString = errMsg
		}

		_, err = s.ChannelMessageSend(m.ChannelID, "```"+pixelString+"```")
		if err != nil {
			errMsg = fmt.Sprintf("could not send message: %v", err)

			fmt.Println(errMsg)
			s.ChannelMessageSend(m.ChannelID, "```"+errMsg+"```")
		}
	}
}
