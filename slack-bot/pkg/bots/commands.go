package bots

import (
	"errors"
	"slack-bot/pkg/config"
	"slack-bot/pkg/utils"

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
		Name: "birth year <year>",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "Calculate your age",
			Examples:    []string{"birth year 1990"},
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				client := botCtx.Client()
				event := botCtx.Event()

				user, err := client.GetUserInfo(event.User)
				if err != nil {
					f := fmt.Sprintf("Failed to get user info: %s\n", err)
					response.ReportError(errors.New(f))
					return
				}

				yearStr := request.Param("year")
				year, err := strconv.Atoi(yearStr)
				if err != nil || year < 0 {
					r := fmt.Sprintf("'%s' is invalid year\n", yearStr)
					response.ReportError(errors.New(r))
				} else {
					age := time.Now().Year() - year
					var r string

					if age < 0 {
						r = fmt.Sprintf("<@%s>You are from the future, go away\n", user.ID)
					} else if age == 0 {
						r = fmt.Sprintf("<@%s> Your Age is %d, You are born this year, really?\n", user.ID, age)
					} else if age < 18 {
						r = fmt.Sprintf("<@%s> Your Age is %d, You are too young\n", user.ID, age)
					} else if age < 22 {
						r = fmt.Sprintf("<@%s> Your Age is %d, So fresh\n", user.ID, age)
					} else if age < 100 {
						r = fmt.Sprintf("<@%s> Your Age is %d, Too old\n", user.ID, age)
					} else {
						r = fmt.Sprintf("<@%s> Your Age is %d, Probably dead\n", user.ID, age)
					}
					response.Reply(r)
				}
			},
		},
	},
	SlackCommand{
		Name: "you suck",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "You can tell the bot that it sucks. But it will talk back.",
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				client := botCtx.Client()
				event := botCtx.Event()

				user, err := client.GetUserInfo(event.User)
				if err != nil {
					f := fmt.Sprintf("Failed to get user info: %s\n", err)
					response.ReportError(errors.New(f))
					return
				}
				r := fmt.Sprintf("<@%s> No, you suck!\nI kwon your IP address btw...", user.ID)
				response.Reply(r)
			},
		},
	},
	SlackCommand{
		Name: "list files",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "List files available for download",
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				fileNames := utils.ListDir(config.ATTACHMENTS_FOLDER)
				if len(fileNames) == 0 {
					response.Reply("No files available for download")
				} else {
					r := fmt.Sprintln("List of files available for download:")
					for _, fileName := range fileNames {
						r += fmt.Sprintf("`%s`\n", fileName)
					}
					response.Reply(r)
				}

			},
		},
	},
}
