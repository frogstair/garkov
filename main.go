package main

import (
	"garkov/bot"
	"os"

	"math/rand"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	rand.Seed(time.Now().Unix())

	godotenv.Load()

	bot.Run(os.Getenv("BOT_KEY"))
}
