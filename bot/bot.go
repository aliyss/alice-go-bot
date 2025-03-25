package bot

import (
	"alice-go-bot/bot/commands"
	"alice-go-bot/bot/commands/bot_commands"
	"alice-go-bot/bot/config"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var discord_bot_commands = map[string]commands.BotCommand{
	bot_commands.EchoCommand.Info.Name:    bot_commands.EchoCommand,
	bot_commands.WeatherCommand.Info.Name: bot_commands.WeatherCommand,
}

func RunBot() {
	// Start your bot here
	bot_config := config.GetConfig()

	bot, err := discordgo.New("Bot " + bot_config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", r.User.Username, r.User.Discriminator)
	})

	// In this example, we only care about receiving message events.
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	regiteredCommands := make(map[string]*commands.BotCommand, len(discord_bot_commands))
	for _, command := range discord_bot_commands {
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, bot_config.GuildID, &command.Info)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", command.Info.Name, err)
		}
		regiteredCommands[cmd.Name] = &command
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := regiteredCommands[i.ApplicationCommandData().Name]; ok {
			h.Handler(s, i)
		}
	})

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
