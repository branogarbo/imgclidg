package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/branogarbo/imgcli/util"
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
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// using janky command code for now

	msgHead := m.Content[:4]

	if msgHead == "!img" {
		imgUrl := m.Content[5:]

		pixelString, err := util.OutputImage(util.OutputConfig{
			Src:          strings.TrimSpace(imgUrl),
			OutputMode:   "ascii",
			AsciiPattern: " .:-=+*#%@",
			IsUseWeb:     true,
			OutputWidth:  66,
		})
		if err != nil {
			errMsg := fmt.Sprintf("imgcli/util: failed to output image, %v", err)

			fmt.Println(errMsg)
			pixelString = errMsg
		}

		_, err = s.ChannelMessageSend(m.ChannelID, "```"+pixelString+"```")
		if err != nil {
			fmt.Println("could not send message,", err)
		}
	}
}
