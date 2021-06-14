package main

import (
	"garkov/bot"
	"garkov/garkov"
	"os"

	"math/rand"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	rand.Seed(time.Now().Unix())

	godotenv.Load()

	go garkov.GarkovLoop()
	bot.Run(os.Getenv("BOT_KEY"))
}
