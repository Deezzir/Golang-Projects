package bots

import (
	"context"
	"slack-bot/pkg/utils"

	"github.com/shomali11/slacker"
)

type Bot interface {
	Init()
	Start(ctx context.Context)

	setCommands(commands []Command)
}

type SlackBot struct {
	Name     string
	BotToken string
	AppToken string

	bot *slacker.Slacker
}

func (s *SlackBot) Init() {
	s.bot = slacker.NewClient(s.BotToken, s.AppToken)
	s.setCommands(SlackCommands...)
}

func (s *SlackBot) logEvents() {
	for event := range s.bot.CommandEvents() {
		utils.CommandLogger.Printf("BOT: %s ", s.Name)
		utils.CommandLogger.Printf("TIME: %s ", event.Timestamp.Format("2000-01-01 12:00:00"))
		utils.CommandLogger.Printf("COMMAND: %s ", event.Command)
		utils.CommandLogger.Printf("PARAMETERS: %s", event.Parameters)
		utils.CommandLogger.Printf("EVENT: %s\n", event.Event)
	}
}

func (s *SlackBot) setCommands(commands ...Command) {
	for _, command := range commands {
		name := command.(SlackCommand).Name
		definition := command.(SlackCommand).CommandDefinition
		s.bot.Command(name, definition)
	}
}

func (s *SlackBot) Start(ctx context.Context) {
	if s.bot != nil {
		go s.logEvents()

		if err := s.bot.Listen(ctx); err != nil {
			panic(err)
		}
	} else {
		panic("SlackBot was not initialized, run Init() first")
	}
}
