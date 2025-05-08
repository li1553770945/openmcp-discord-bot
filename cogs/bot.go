package cogs

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

var (
	botOnce sync.Once
)

func GetGlobalBot() {
	botOnce.Do(func() {
		discord, err := discordgo.New("Bot " + "authentication token")
	})
}
