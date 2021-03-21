// package main

// import (
// 	"flag"
// 	"log"

// 	dg "github.com/bwmarrin/discordgo"
// )

// var (
// 	token   string
// 	session *dg.Session
// 	err     error
// )

// func init() {
// 	flag.StringVar(&token, "t", "", "Bot access token.")
// 	flag.Parse()
// }

// func main() {
// 	session, err = dg.New("Bot " + token)
// 	if err != nil {
// 		log.Fatal("failed to start bot:", err)
// 	}

// 	session.Open()
// 	defer session.Close()

// 	session.AddHandler(func(s *dg.Session, e *dg.MessageCreate) {
// 		var (
// 			msg string = e.Message.Content
// 		)

// 		_, err = s.ChannelMessageSend("imgcli", msg)
// 		if err != nil {
// 			log.Fatal("failed to send message:", err)
// 		}

// 	})
// }

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
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

	pixelString, err := util.OutputImage(util.OutputConfig{
		Src:          "https://static.wikia.nocookie.net/among-us-wiki/images/8/84/Among_Us.png/revision/latest?cb=20201019142953",
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

	if m.Content == "!img" {
		s.ChannelMessageSend(m.ChannelID, "```"+pixelString+"```")
	}
}
