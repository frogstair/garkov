package bot

import (
	"fmt"
	"garkov/garkov"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

var PREFIX = "--"

func Run(token string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(onMessage)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot started")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
	fmt.Println("\nStopping...")
	dg.Close()
	fmt.Println("Bot stopped")
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, PREFIX) {
		return
	}

	s.ChannelTyping(m.ChannelID)

	message := strings.Trim(m.Content, PREFIX+" ")
	space := regexp.MustCompile(`\s+`)
	message = space.ReplaceAllString(message, " ")
	args := strings.Split(message, " ")

	if args[0] == "prefix" {
		PREFIX = args[1]
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set prefix to %s", PREFIX))
		return
	}

	if args[0] == "garkov" {
		path := garkov.Garkov()
		file, err := os.Open("temp/" + path)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Internal error")
			panic(err)
		}
		s.ChannelFileSend(m.ChannelID, fmt.Sprintf("garkov%d.png", time.Now().Unix()), file)
		os.Remove("temp/" + path)
		return
	}

	if args[0] == "help" {
		e := embed.NewGenericEmbedAdvanced("Help", "These are the commands\nhelp: show this message\ngarkov: generate an image\nprefix: change prefix. Current is "+PREFIX, 0xa4781c)
		s.ChannelMessageSendEmbed(m.ChannelID, e)
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Unknown command "+fmt.Sprint(args))

}
