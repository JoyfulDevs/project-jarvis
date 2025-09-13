package app

import (
	"context"

	"github.com/joyfuldevs/project-jarvis/pkg/slack"
	"github.com/joyfuldevs/project-jarvis/pkg/slack/bot"
)

var _ bot.EventHandler = (*JarvisBot)(nil)

type JarvisBot struct {
	AppToken string
	BotToken string
}

func NewJarvisBot(appToken, botToken string) *JarvisBot {
	return &JarvisBot{
		AppToken: appToken,
		BotToken: botToken,
	}
}

func (j *JarvisBot) Run(ctx context.Context) error {
	bot := bot.NewBot(j.AppToken, j.BotToken, j)
	return bot.Run(ctx)
}

func (j *JarvisBot) HandleCommandEvent(payload *slack.SlashCommandEventPayload) {
	responder := &CommandResponder{
		AppToken: j.AppToken,
		BotToken: j.BotToken,
		Payload:  payload,
	}
	responder.RespondCommand()
}

func (j *JarvisBot) HandleInteractiveEvent(payload *slack.InteractiveEventPayload) {
	responder := &ActionResponder{
		AppToken: j.AppToken,
		BotToken: j.BotToken,
		Payload:  payload,
	}
	responder.RespondActions()
}
