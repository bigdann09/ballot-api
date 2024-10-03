package bot

import (
	"fmt"

	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"
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
	username := message.Chat.UserName
	premium := message.From.IsPremium

	// check user
	if !models.CheckUser(chatID) {
		_, err := models.NewUser(&models.User{
			TGID:      chatID,
			TGPremium: premium,
			Token:     utils.ReferralToken(8),
		})

		if err != nil {
			panic(err)
		}

		msg := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("./assets/cover.jpg"))
		msg.Caption = fmt.Sprintf("Welcome %s 😸 to the most engaging miniapp platform on telegram", username)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyMarkup = WelcomeNewUserMarkup()
		bot.api.Send(msg)

		return
	}

	msg := tgbotapi.NewMessage(chatID, "hello")
	bot.api.Send(msg)
}
