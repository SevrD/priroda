package core

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/domain"
	"main/internal/models"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Core struct {
	users        domain.Users
	chat         domain.ChatStatus
	announcement domain.Announcement
	bot          *telego.Bot
	chatIDAdmin  int64
}

func NewCore(users domain.Users, chat domain.ChatStatus, announcement domain.Announcement, bot *telego.Bot) domain.Core {
	return &Core{
		users:        users,
		chat:         chat,
		announcement: announcement,
		bot:          bot,
	}
}

func (c *Core) SendMessageAfterAddAnnouncement(ctx context.Context, tgUserID int64, chatID telego.ChatID, fileID *string) {

	if c.chatIDAdmin == 0 {
		return
	}

	annID := c.chat.GetAnnId(ctx, tgUserID)

	text, _, err := c.announcement.GetAnnouncement(ctx, tgUserID, annID)
	if err != nil {
		log.Println("Get text announcements from BD:", err)
		c.users.SendError(chatID)
		return
	}

	text, err = c.addContacts(ctx, tgUserID, text)

	if err != nil {
		log.Println("Get user info error:", err)
		c.users.SendError(chatID)
		return
	}

	var inlineKeyboardRow []telego.InlineKeyboardButton

	inlineKeyboardRow = append(inlineKeyboardRow, tu.InlineKeyboardButton("Удалить").WithCallbackData("delete_ann"))
	inlineKeyboardRow = append(inlineKeyboardRow, tu.InlineKeyboardButton("Бан").WithCallbackData("ban"))
	inlineKeyboardRow = append(inlineKeyboardRow, tu.InlineKeyboardButton("Снять бан").WithCallbackData("unban"))
	inlineKeyboard := tu.InlineKeyboard(tu.InlineKeyboardRow(inlineKeyboardRow...))

	chatIDAdmin := tu.ID(c.chatIDAdmin)

	var adm_msq *telego.Message

	if fileID != nil {

		err = c.announcement.AddPhoto(ctx, annID, *fileID)
		if err != nil {
			log.Println("Save photo:", err)
			c.users.SendError(chatID)
			return
		}

		photoFile := tu.FileFromID(*fileID)

		photoParams := tu.Photo(chatIDAdmin, photoFile).WithCaption(text).WithReplyMarkup(inlineKeyboard)

		adm_msq, err = c.bot.SendPhoto(photoParams)

	} else {
		message := tu.Message(chatIDAdmin, text).WithReplyMarkup(inlineKeyboard)
		adm_msq, err = c.bot.SendMessage(message)
	}

	if err != nil {
		log.Println("Send announcement:", err)
		c.users.SendError(chatID)
		return
	}

	err = c.announcement.SetAdminMsgID(ctx, annID, int64(adm_msq.MessageID))

	if err != nil {
		log.Println("Save admin message id:", err)
		c.users.SendError(chatID)
		return
	}

	c.publicAnnouncement(ctx, adm_msq)

	msg := fmt.Sprintf("Номер вашего объявления: %d", annID)
	newMessage := tu.Message(chatID, msg)
	c.bot.SendMessage(newMessage)

	c.chat.Save(ctx, tgUserID, models.StatusCode(0), 0)

}

func (c *Core) addContacts(ctx context.Context, tgUserID int64, txt string) (string, error) {

	name, login, _, err := c.users.GetUserInfo(ctx, tgUserID)

	if err != nil {
		return "", err
	}

	if name != "" {
		txt = txt + "\nАвтор: " + name
	}

	if login != "" {
		txt = txt + "\nЛогин: @" + login
	}

	return txt, nil
}

func (c *Core) publicAnnouncement(ctx context.Context, Message telego.MaybeInaccessibleMessage) {

	channelID := tu.ID(config.AppConfig.ChannelID)

	chatID := tu.ID(Message.GetChat().ID)

	annInfo, err := c.announcement.GetAnnouncementOnAdmMsgID(ctx, int64(Message.GetMessageID()))

	if err != nil {
		log.Println("Get announcement error:", err)
		c.users.SendError(chatID)
		return
	}

	text, err := c.addContacts(ctx, annInfo.TgID, annInfo.Text)

	if err != nil {
		log.Println("Get user info error:", err)
		c.users.SendError(chatID)
		return
	}

	if annInfo.FileID != "" {
		photo := tu.FileFromID(annInfo.FileID)
		photoParams := tu.Photo(channelID, photo).WithCaption(text)
		msg, err := c.bot.SendPhoto(photoParams)
		if err != nil {
			log.Println("Public announcement error:", err)
			c.users.SendError(chatID)
			return
		}
		c.announcement.SetPublicID(ctx, annInfo.Id, int64(msg.MessageID))
	} else {
		message := tu.Message(channelID, text)
		pMessage, err := c.bot.SendMessage(message)
		if err != nil {
			log.Println("Public announcement error:", err)
			c.users.SendError(chatID)
			return
		}

		c.announcement.SetPublicID(ctx, annInfo.Id, int64(pMessage.GetMessageID()))
	}

}

func (c *Core) SendMessageWithButtons(chatID telego.ChatID, text string, buttons []models.Button) {

	var inlineKeyboardRow []telego.InlineKeyboardButton

	for _, value := range buttons {
		inlineKeyboardRow = append(inlineKeyboardRow, tu.InlineKeyboardButton(value.Text).WithCallbackData(value.Name))
	}

	// Inline keyboard parameters
	inlineKeyboard := tu.InlineKeyboard(tu.InlineKeyboardRow(inlineKeyboardRow...))

	// Message parameters
	message := tu.Message(
		chatID,
		text,
	).WithReplyMarkup(inlineKeyboard)

	// Sending message
	c.bot.SendMessage(message)

}

func (c *Core) SendDefaultMessage(chatID telego.ChatID) {

	text := `Воспользуйтесь кнопками из меню для добавления и удаления объявлений или используйте команды:
/add Добавить новое объявление
/delete Удалить объявление`

	var buttons []models.Button
	buttons = append(buttons, models.Button{Text: "Добавить объявление", Name: "add"})
	buttons = append(buttons, models.Button{Text: "Удалить объявление", Name: "delete"})

	c.SendMessageWithButtons(chatID, text, buttons)
}

func (c *Core) SaveAdminChatID(chatID int64) {
	c.chatIDAdmin = chatID
}
