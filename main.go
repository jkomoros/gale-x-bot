package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

const APP_NAME = "gale-x-bot"
const TOKEN_ENV_NAME = "BOT_TOKEN"

//The name of the category this will look for that contains things this should treat as threads
const THREAD_CATEGORY_NAME = "Threads"
const THREAD_ARCHIVE_CATEGORY_NAME = THREAD_CATEGORY_NAME + " Archive"

//Number of threads to allow in an active category before bumping old ones.
//Should be at least one smaller than MAX_CATEGORY_CHANNELS.
const MAX_ACTIVE_THREADS = 5

//The max number of channels that Discord allows to be in a category channel. This is configured by discord.
const MAX_CATEGORY_CHANNELS = 50

var token string

func main() {

	if token == "" {
		token = os.Getenv(TOKEN_ENV_NAME)
	}

	if token == "" {
		fmt.Println("No token provided. Please run: " + APP_NAME + " -t <bot token> or set env var " + TOKEN_ENV_NAME)
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	newBot(dg)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println(APP_NAME + " is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
