package bot

import (
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

func (bot *Bot) handleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		bot.handleCommand(message)
	}
}

func (bot *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		bot.handleStartCommand(message)
	}
}

func (bot *Bot) handleStartCommand(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Know your stand, let's get started and decide who will rule over the world's most influcencial nation")
	msg.ReplyMarkup = StartUserMarkup()
	bot.api.Send(msg)
}
