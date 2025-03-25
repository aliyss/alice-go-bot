package bot_commands

import (
	"alice-go-bot/bot/commands"

	"github.com/bwmarrin/discordgo"
)

func echoCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: i.ApplicationCommandData().Options[0].StringValue(),
		},
	})
}

var EchoCommand commands.BotCommand = commands.BotCommand{
	Info: discordgo.ApplicationCommand{
		Name:        "echo",
		Description: "Replies with your message!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message to echo back",
				Required:    true,
			},
		},
	},
	Handler: echoCommandHandler,
}
