package bot

import (
	"github.com/ballot/internals/utils"
	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

func StartUserMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Join Channel", utils.CHANNEL_URL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Join Group", utils.GROUP_URL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Follow On X", utils.TWITTER_URL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp("Let's go", tgbotapi.WebAppInfo{
				URL: utils.MINIAPP_URL,
			}),
		),
	)
}
