package bots

import (
	"github.com/shomali11/slacker"

	"fmt"
	"strconv"
	"time"
)

type Command interface {
}

type SlackCommand struct {
	Name              string
	CommandDefinition *slacker.CommandDefinition
}

var SlackCommands = []Command{
	SlackCommand{
		Name: "ping",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "Ping the bot",
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				response.Reply("Pong!")
			},
		},
	},

	SlackCommand{
		Name: "birth year <year>",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "Calculate your age",
			Examples:    []string{"birth year 1990"},
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				yearStr := request.Param("year")
				year, err := strconv.Atoi(yearStr)
				if err != nil {
					r := fmt.Sprintf("'%s' is not a year\n", yearStr)
					response.Reply(r)
				} else {
					age := time.Now().Year() - year
					var r string

					if age > 21 {
						r = fmt.Sprintf("Your Age is %d, Too old\n", age)
					} else {
						r = fmt.Sprintf("Your Age is %d, So fresh\n", age)
					}
					response.Reply(r)
				}
			},
		},
	},

	SlackCommand{
		Name:              "help",
		CommandDefinition: &slacker.CommandDefinition{},
	},

	SlackCommand{
		Name: "you suck",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "You can tell the bot that it sucks. But it will talk back.",
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				response.Reply("You suck!")
			},
		},
	},
}
