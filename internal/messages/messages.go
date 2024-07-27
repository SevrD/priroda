package messages

import (
	"context"
	"log"
	"main/internal/config"
	"main/internal/domain"
	"main/internal/models"
	"strconv"

	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/mymmrac/telego"
)

type Messages struct {
	announcement domain.Announcement
	users        domain.Users
	chat         domain.ChatStatus
	commands     domain.CommandsWorker
	core         domain.Core
	bot          *telego.Bot
}

func NewMessages(announcement domain.Announcement, users domain.Users, chat domain.ChatStatus, commands domain.CommandsWorker, core domain.Core, bot *telego.Bot) *Messages {
	return &Messages{
		announcement: announcement,
		users:        users,
		chat:         chat,
		commands:     commands,
		core:         core,
		bot:          bot,
	}
}

func (m *Messages) Text(ctx context.Context, message *telego.Message) {

	chatID := tu.ID(message.Chat.ID)
	var ann_id int64
	var err error
	if len(message.Photo) > 0 && message.Caption != "" {
		ann_id, err = m.announcement.Add(ctx, message.From.ID, message.Caption, message.Chat.ID)

	} else if message.Text != "" {
		ann_id, err = m.announcement.Add(ctx, message.From.ID, message.Text, message.Chat.ID)
	} else {
		return
	}
	if err != nil {
		log.Println("Announcement add error:", err)
		m.users.SendError(message.Chat.ChatID())
		return
	}
	err = m.chat.Save(ctx, message.From.ID, 2, ann_id) // запрос картинок
	if err != nil {
		log.Println("Status save error:", err)
		m.users.SendError(message.Chat.ChatID())
		return
	}

	var buttons []models.Button
	buttons = append(buttons, models.Button{Text: "Пропустить", Name: "skip"})

	m.core.SendMessageWithButtons(chatID, "Отправьте ОДНУ фотографию или пропустите этот шаг. Остальные фото можно добавить после публикации объявления.", buttons)

}

func (m *Messages) Photo(ctx context.Context, message *telego.Message) {

	chatID := tu.ID(message.Chat.ID)

	if message.Photo != nil {
		fileID := message.Photo[len(message.Photo)-1].FileID
		m.core.SendMessageAfterAddAnnouncement(ctx, message.From.ID, chatID, &fileID)
		m.chat.Save(ctx, message.From.ID, models.StatusCode(0), 0)

	} else {
		if message.Document != nil {
			message := tu.Message(chatID, "Файлы не принимаются. Отправьте фото")
			m.bot.SendMessage(message)
		}
		var buttons []models.Button
		buttons = append(buttons, models.Button{Text: "Пропустить", Name: "skip"})
		m.core.SendMessageWithButtons(chatID, "Отправьте ОДНУ фотографию или пропустите этот шаг. Остальные фото можно добавить после публикации объявления.", buttons)
	}

}

func (m *Messages) Number(ctx context.Context, message *telego.Message) {

	chatID := tu.ID(message.Chat.ID)

	id, err := strconv.Atoi(message.Text)
	if err != nil {
		message := tu.Message(chatID, "Необходимо ввести целое число")
		m.bot.SendMessage(message)
	}

	_, msgID, err := m.announcement.GetAnnouncement(ctx, message.From.ID, int64(id))

	if err != nil {
		message := tu.Message(
			chatID,
			"Объявление с таким номером не найдено или оно другого пользователя",
		)
		m.bot.SendMessage(message)
		return
	}

	channelID := tu.ID(config.AppConfig.ChannelID)
	deleteMessageParams := &telego.DeleteMessageParams{ChatID: channelID, MessageID: int(msgID)}
	err = m.bot.DeleteMessage(deleteMessageParams)

	if err != nil {
		log.Println("Delete announcement error:", err)
		m.users.SendError(message.Chat.ChatID())
	} else {
		message := tu.Message(
			chatID,
			"Объявление удалено",
		)
		m.bot.SendMessage(message)
	}

	m.chat.Save(ctx, message.From.ID, models.StatusCode(0), 0)
}
