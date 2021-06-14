package bot

import (
	"fmt"
	"garkov/garkov"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

var prefices map[string]string

var DPREFIX = "--"

func Run(token string) {
	prefices = make(map[string]string)
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

	prefix, ok := prefices[m.GuildID]
	if !ok {
		prefix = DPREFIX
	}

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	s.ChannelTyping(m.ChannelID)

	message := strings.Trim(m.Content, prefix+" ")
	space := regexp.MustCompile(`\s+`)
	message = space.ReplaceAllString(message, " ")
	args := strings.Split(message, " ")

	if args[0] == "prefix" {
		if len(args) == 1 {
			s.ChannelMessageSend(m.ChannelID, "You need to specify a prefix (e.g. "+prefix+"prefix ?)")
			return
		}

		prefices[m.GuildID] = args[1]
		prefix = args[1]
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set prefix to %s", prefix))
		return
	}

	if args[0] == "garkov" {
		path := nextGarkov()
		file, err := os.Open("cache/" + path)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Internal error")
			panic(err)
		}
		s.ChannelFileSend(m.ChannelID, fmt.Sprintf("garkov%d.png", time.Now().Unix()), file)
		os.Remove("cache/" + path)
		return
	}

	if args[0] == "help" {
		e := embed.NewGenericEmbedAdvanced("Help", "These are the commands\n__**help**__: show this message\n__**garkov**__: generate an image\n__**prefix**__: change prefix. Current is "+prefix+"\n__**donate**__: send a donation to me!", 0xa4781c)
		s.ChannelMessageSendEmbed(m.ChannelID, e)
		return
	}

	if args[0] == "donate" {
		e := new(discordgo.MessageEmbed)
		e.Title = "Help support Garkov!"
		e.Color = 0xa4781c
		i := new(discordgo.MessageEmbedImage)
		i.URL = "https://i.imgur.com/Fe006FE.png"
		e.Image = i
		e.Description = "Send a **donation** to my [PayPal](https://paypal.me/frogstair) or scan the QR code"
		s.ChannelMessageSendEmbed(m.ChannelID, e)
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Unknown command "+fmt.Sprint(args))

}

func nextGarkov() string {
	for len(garkov.ImageChannel) == 0 {
		log.Printf("Queue is empty, waiting")
		time.Sleep(500 * time.Millisecond)
	}
	v := <-garkov.ImageChannel
	log.Printf("Removed item %s", v)
	return v
}
