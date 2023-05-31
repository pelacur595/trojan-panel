package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
	"trojan-panel/core"
)

var bot = new(tgbotapi.BotAPI)

func InitTelegramBotApi() {
	var err error
	config := core.Config
	bot, err = tgbotapi.NewBotAPI(os.Getenv(config.ApiToken))
	if err != nil {
		logrus.Errorf("new bot api err: %v", err)
		panic(err)
	}
	logrus.Infof("Authorized on account %s", bot.Self.UserName)
}

func GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return bot.GetUpdatesChan(u)
}
