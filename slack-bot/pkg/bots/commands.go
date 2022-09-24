package bots

import (
	"slack-bot/pkg/config"
	"slack-bot/pkg/utils"
	"strings"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"

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
					utils.ErrorLogger.Printf("Failed to get user info: %s\n", err)
					response.Reply("Something went wrong, sorry")
					return
				}

				yearStr := request.Param("year")
				year, err := strconv.Atoi(yearStr)
				if err != nil || year < 0 {
					r := fmt.Sprintf("'%s' is an invalid year\n", yearStr)
					response.Reply(r)
				} else {
					age := time.Now().Year() - year
					var r string

					if age < 0 {
						r = fmt.Sprintf("<@%s>You are from the future, go away\n", user.ID)
					} else if age == 0 {
						r = fmt.Sprintf("<@%s> Your Age is *%d*, You are born this year, really?\n", user.ID, age)
					} else if age < 18 {
						r = fmt.Sprintf("<@%s> Your Age is *%d*, You are too young\n", user.ID, age)
					} else if age < 22 {
						r = fmt.Sprintf("<@%s> Your Age is *%d*, So fresh\n", user.ID, age)
					} else if age < 100 {
						r = fmt.Sprintf("<@%s> Your Age is *%d*, Too old\n", user.ID, age)
					} else {
						r = fmt.Sprintf("<@%s> Your Age is *%d*, Probably dead\n", user.ID, age)
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
					utils.ErrorLogger.Printf("Failed to get user info: %s\n", err)
					response.Reply("Something went wrong, sorry")
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
					return
				}

				r := fmt.Sprintln("List of files available for download:")
				for _, fileName := range fileNames {
					r += fmt.Sprintf("`%s`\n", fileName)
				}
				response.Reply(r)
			},
		},
	},
	SlackCommand{
		Name: "get file <file>",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "Get available file",
			Examples:    []string{"get file dog.jpg", "get file doc.pdf"},
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				file := request.Param("file")
				client := botCtx.Client()
				event := botCtx.Event()

				if filePath, ok := utils.GetFile(config.ATTACHMENTS_FOLDER, file); ok {
					if event.Channel != "" {
						params := slack.FileUploadParameters{
							Content:  filePath,
							Channels: []string{event.Channel},
						}

						client.PostMessage(event.Channel, slack.MsgOptionText("Downloading file ...", false))
						_, err := client.UploadFile(params)
						if err != nil {
							utils.ErrorLogger.Printf("Failed to upload '%s' file to Slack channel\n", filePath)
							response.Reply("Sorry, failed to download the file :'(")
						}
					}
				} else {
					response.Reply("File not found, use `list files` for avaliable files")
				}
			},
		},
	},
	SlackCommand{
		Name: "put file <filename> <description>",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "Saves the provided file to Noxu-bot's memory",
			Examples:    []string{"put file dog.jpeg"},
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {

			},
		},
	},
	SlackCommand{
		Name: "validate <email>",
		CommandDefinition: &slacker.CommandDefinition{
			Description: "Check that email is valid and verifies email domain. Does not check if email exists",
			Examples:    []string{"validate deezzir@gmail.com"},
			Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
				param := request.Param("email")
				email := utils.GetHyperlinkTxt(param)
				if email == "" {
					email = param
				}

				local, domain, ok := utils.NormalizeEmail(email)
				if !ok {
					response.Reply("Please provide a valid email")
					return
				}
				r := fmt.Sprintf("*Email*: `%s@%s`\n", local, domain)
				r += fmt.Sprintf("*Domain*: `%s`\n", domain)

				vdDomain := utils.CheckEmailDomain(domain)

				if vdDomain.Valid {
					if len(vdDomain.Addr) > 0 {
						r += fmt.Sprintf("- *Addresses*: `%s`\n", strings.Join(vdDomain.Addr[:], "`, `"))
					}
					r += fmt.Sprintf("- *has MX*: `%t`\n", vdDomain.HasMX)

					if vdDomain.HasSPF {
						r += fmt.Sprintf("- *SPF Record*: `%s`\n", vdDomain.SPFRecord)
					}

					if vdDomain.HasDMARC {
						r += fmt.Sprintf("- *DMARC Record*: `%s`\n", vdDomain.DMARCRecord)
					}

				} else {
					r += fmt.Sprintf("- *Valid*: `%t`\n", vdDomain.Valid)
				}

				response.Reply(r)
			},
		},
	},
}
