package messages

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/domain"
	"main/internal/models"
	"strconv"
	"unicode/utf8"

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
	var fullText string
	if len(message.Photo) > 0 && message.Caption != "" {
		ann_id, err = m.announcement.Add(ctx, message.From.ID, message.Caption, message.Chat.ID)
	} else if message.Text != "" {
		fullText, err = m.core.AddContacts(ctx, message.From.ID, message.Text)
		if err != nil {
			log.Println("Add contacts to text error:", err)
			m.users.SendError(message.Chat.ChatID())
			return
		}
		lenCnt := utf8.RuneCountInString(m.core.Contacts(ctx, message.From.ID))
		lenTxt := utf8.RuneCountInString(fullText)
		if lenTxt > 1024-lenCnt-1 {
			bidText(m, lenTxt-(1024-lenCnt-1), chatID)
			return
		}
		ann_id, err = m.announcement.Add(ctx, message.From.ID, message.Text, message.Chat.ID)
	} else {
		return
	}
	if err != nil {
		log.Println("Add announcement error:", err)
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
		msg := tu.Message(chatID, "Необходимо ввести целое число")
		m.bot.SendMessage(msg)
	}

	_, msgID, _, err := m.announcement.GetAnnouncement(ctx, message.From.ID, int64(id))

	if err != nil {
		msg := tu.Message(
			chatID,
			"Объявление с таким номером не найдено или оно другого пользователя",
		)
		m.bot.SendMessage(msg)
		return
	}

	channelID := tu.ID(config.AppConfig.ChannelID)
	deleteMessageParams := &telego.DeleteMessageParams{ChatID: channelID, MessageID: int(msgID)}
	err = m.bot.DeleteMessage(deleteMessageParams)

	if err != nil {
		if fmt.Sprint(err) == "telego: deleteMessage(): api: 400 \"Bad Request: message to delete not found\"" || fmt.Sprint(err) == "telego: deleteMessage(): api: 400 \"Bad Request: message can`t be deleted\"" {
			err = m.core.SendDeleteRequest(ctx, message.From.ID, int64(id), chatID)
			if err != nil {
				log.Println("Send delete request error:", err)
				m.users.SendError(message.Chat.ChatID())
			}
			msg := tu.Message(chatID, "Прошло больше 48 часов с момента публикации. Запрос на удаление отправлен @"+m.users.LoginAdmin())
			m.bot.SendMessage(msg)
		} else {
			log.Println("Delete announcement error:", err)
			m.users.SendError(message.Chat.ChatID())
		}

	} else {
		message := tu.Message(chatID, "Объявление удалено")
		m.bot.SendMessage(message)
	}

	m.chat.Save(ctx, message.From.ID, models.StatusCode(0), 0)
}

func declOfNum(number int, titles []string) string {
	cases := []int{2, 0, 1, 1, 1, 2}
	var currentCase int
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	return titles[currentCase]
}

func bidText(m *Messages, len int, chatID telego.ChatID) {
	titles := []string{"символ", "символа", "символов"}
	text := "Размер текста не должен превышать 1024 символа. Необходимо уменьшить текст на %d %s и отправить в следующем сообщении. Более подробную информацию можно указать в комментариях опубликованного объявления."
	text = fmt.Sprintf(text, len, declOfNum(len, titles))
	message := tu.Message(chatID, text)
	m.bot.SendMessage(message)
}
