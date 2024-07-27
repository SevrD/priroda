package commands

import (
	"context"
	"log"
	"main/internal/config"
	"main/internal/domain"
	"main/internal/models"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Commands struct {
	bot   *telego.Bot
	users domain.Users
	chat  domain.ChatStatus
	core  domain.Core
}

func NewCommands(bot *telego.Bot, users domain.Users, chat domain.ChatStatus, core domain.Core) domain.CommandsWorker {
	return &Commands{
		bot:   bot,
		users: users,
		chat:  chat,
		core:  core,
	}
}

func (c *Commands) Start(ctx context.Context, chatID telego.ChatID, tgID int64, rules string) {

	c.chat.Save(ctx, tgID, models.StatusCode(0), 0)

	var buttons []models.Button
	buttons = append(buttons, models.Button{Text: "Принимаю", Name: "accept"})

	c.core.SendMessageWithButtons(chatID, "Ознакомьтесь с правилами подачи объявления и примите их."+rules, buttons)

}

func (c *Commands) Add(ctx context.Context, chatID telego.ChatID, tgID int64) {

	c.chat.Save(ctx, tgID, models.StatusCode(0), 0)

	_, _, ban, err := c.users.GetUserInfo(ctx, tgID)

	if err != nil {
		message := tu.Message(
			chatID,
			"Сначала необходимо принять правила публикации объявлений. /start",
		)

		c.bot.SendMessage(message)
		return
	}

	if ban {
		message := tu.Message(
			chatID,
			"Вы заблокированы. Соблюдайте правила /start. Для разблокировки обратитесь к администратору @"+config.AppConfig.LoginAdmin,
		)

		c.bot.SendMessage(message)
		return
	}

	err = c.chat.Save(ctx, tgID, models.StatusCode(1), 0) // запрос текста объявления
	if err != nil {
		log.Println("Register error:", err)
		c.users.SendError(chatID)
		return
	}
	message := tu.Message(
		chatID,
		"Напишите текст объявления одним сообщением без картинок. Не забудьте указать цену и контактную информацию.",
	)

	c.bot.SendMessage(message)
}

func (c *Commands) Delete(ctx context.Context, chatID telego.ChatID, tgID int64) {

	c.chat.Save(ctx, tgID, models.StatusCode(0), 0)

	message := tu.Message(
		chatID,
		"Отправьте номер удаляемого объявления числом в следующем сообщении",
	)
	c.bot.SendMessage(message)

	err := c.chat.Save(ctx, tgID, models.StatusCode(3), 0) // запрос текста объявления
	if err != nil {
		log.Println("Save status error:", err)
		c.users.SendError(chatID)
		return
	}
}
