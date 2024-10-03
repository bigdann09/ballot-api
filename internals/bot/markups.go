package bot

import (
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

func WelcomeNewUserMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp("Launch app 🔫", tgbotapi.WebAppInfo{
				URL: "https://www.minicatbook.com",
			}),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Join telegram community", "https://t.me/catbook"),
		),
	)
}
