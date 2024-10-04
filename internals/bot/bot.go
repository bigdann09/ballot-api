package bot

import (
	"github.com/ballot/internals/utils"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
}

func New(token string) (*Bot, error) {
	tgbot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &Bot{}, err
	}

	tgbot.Send(tgbotapi.SetChatMenuButtonConfig{
		MenuButton: &tgbotapi.MenuButton{
			Type: "web_app",
			Text: "Vote",
			WebApp: &tgbotapi.WebAppInfo{
				URL: utils.MINIAPP_URL,
			},
		},
	})

	return &Bot{api: tgbot}, nil
}

func (bot *Bot) HandleUpdates(update *tgbotapi.Update) {
	switch {
	case update.Message != nil:
		// handle all message request
		bot.handleMessage(update.Message)
	}
}
