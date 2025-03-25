package commands

import (
	"github.com/bwmarrin/discordgo"
)

type BotCommand struct {
	Info    discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
