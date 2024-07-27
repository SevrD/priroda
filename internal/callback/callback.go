package callback

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"
	"main/internal/domain"
	"main/internal/models"
	"main/internal/queries"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Callback struct {
	queries      *queries.Queries
	announcement domain.Announcement
	users        domain.Users
	bot          *telego.Bot
	core         domain.Core
	chat         domain.ChatStatus
	commands     domain.CommandsWorker
}

func NewCallBack(queries *queries.Queries, announcement domain.Announcement, users domain.Users, bot *telego.Bot, core domain.Core, chat domain.ChatStatus, commands domain.CommandsWorker) domain.CallBacker {
	return &Callback{
		queries:      queries,
		announcement: announcement,
		users:        users,
		bot:          bot,
		core:         core,
		chat:         chat,
		commands:     commands,
	}
}

func (c *Callback) CallBack(ctx context.Context, callbackQuery *telego.CallbackQuery) {

	chatID := tu.ID(callbackQuery.From.ID)

	switch callbackQuery.Data {
	case "accept": // приняты правила
		c.register(ctx, callbackQuery)
	case "ban":
		c.ban(ctx, callbackQuery)
	case "unban":
		c.unban(ctx, callbackQuery)
	case "skip": // пропущена картинка
		c.core.SendMessageAfterAddAnnouncement(ctx, callbackQuery.From.ID, chatID, nil)
	case "delete_ann": // удалено объявление админом
		c.deleteAnnFromAdmin(ctx, callbackQuery)
	case "add": // добавить объявление
		c.commands.Add(ctx, chatID, callbackQuery.From.ID)
	case "delete":
		c.commands.Delete(ctx, chatID, callbackQuery.From.ID)
	}

}

func addReaction(bot *telego.Bot, Message telego.MaybeInaccessibleMessage, emoji string) {
	var reaction []telego.ReactionType
	reaction = append(reaction, &telego.ReactionTypeEmoji{Type: "emoji", Emoji: emoji})

	SetMessageReactionParams := &telego.SetMessageReactionParams{
		ChatID:    tu.ID(Message.GetChat().ID),
		MessageID: Message.GetMessageID(),
		Reaction:  reaction,
	}
	bot.SetMessageReaction(SetMessageReactionParams)
}

func (c *Callback) register(ctx context.Context, callbackQuery *telego.CallbackQuery) {

	chatID := tu.ID(callbackQuery.From.ID)

	err := c.users.Register(ctx,
		callbackQuery.From.ID,
		callbackQuery.From.Username,
		callbackQuery.From.FirstName,
		int(callbackQuery.Message.GetDate()),
		callbackQuery.Message.GetChat().ID)
	if err != nil {
		log.Println("Register error:", err)
		c.users.SendError(chatID)
		return
	}

	myCommandsParams := &telego.GetMyCommandsParams{} //Scope: botCommandScope}
	cm, err := c.bot.GetMyCommands(myCommandsParams)

	if err != nil {
		log.Println("Get commands error:", err)
		c.users.SendError(chatID)
		return
	}

	cm = append(cm, telego.BotCommand{Command: "addAnnouncement", Description: "Добавить объявление"})
	cm = append(cm, telego.BotCommand{Command: "delAnnouncement", Description: "Удалить объявление"})

	SetMyCommandsParams := &telego.SetMyCommandsParams{Commands: cm, Scope: &telego.BotCommandScopeAllGroupChats{Type: "all_group_chats"}} //, Scope: botCommandScope}
	c.bot.SetMyCommands(SetMyCommandsParams)

	c.core.SendDefaultMessage(chatID)

	c.chat.Save(ctx, callbackQuery.From.ID, models.StatusCode(0), 0)
}

func (c *Callback) ban(ctx context.Context, callbackQuery *telego.CallbackQuery) {

	chatID := tu.ID(callbackQuery.From.ID)
	annInfo, err := c.announcement.GetAnnouncementOnAdmMsgID(ctx, int64(callbackQuery.Message.GetMessageID()))

	if err != nil {
		log.Println("Get announcement info error:", err)
		c.users.SendError(chatID)
		return
	}

	err = c.users.Ban(ctx, annInfo.TgID)

	if err != nil {
		log.Println("Ban user error:", err)
		c.users.SendError(chatID)
		return
	}

	addReaction(c.bot, callbackQuery.Message, "💩")

	chatID = tu.ID(annInfo.ChatID)

	msg := "На Вас наложен бан.\nВы больше не можете добавлять объявления.\nСмотрите правила публикации /start"
	message := tu.Message(chatID, msg)
	c.bot.SendMessage(message)
}

func (c *Callback) unban(ctx context.Context, callbackQuery *telego.CallbackQuery) {

	chatID := tu.ID(callbackQuery.From.ID)
	annInfo, err := c.announcement.GetAnnouncementOnAdmMsgID(ctx, int64(callbackQuery.Message.GetMessageID()))

	if err != nil {
		log.Println("Get announcement info error:", err)
		c.users.SendError(chatID)
		return
	}

	err = c.users.UnBan(ctx, annInfo.TgID)

	if err != nil {
		log.Println("Unban user error:", err)
		c.users.SendError(chatID)
		return
	}

	addReaction(c.bot, callbackQuery.Message, "🤝")

	chatID = tu.ID(annInfo.ChatID)

	msg := "Бан снят"
	message := tu.Message(chatID, msg)
	c.bot.SendMessage(message)
}

func (c *Callback) deleteAnnFromAdmin(ctx context.Context, callbackQuery *telego.CallbackQuery) {

	chatID := tu.ID(callbackQuery.From.ID)

	annInfo, err := c.announcement.GetAnnouncementOnAdmMsgID(ctx, int64(callbackQuery.Message.GetMessageID()))

	if err != nil {
		log.Println("Get announcement info error:", err)
		c.users.SendError(chatID)
		return
	}

	channelID := tu.ID(config.AppConfig.ChannelID)

	deleteMessageParams := &telego.DeleteMessageParams{ChatID: channelID, MessageID: int(annInfo.PublicID)}
	err = c.bot.DeleteMessage(deleteMessageParams)

	if err != nil {
		log.Println("Delete announcement error:", err)
		c.users.SendError(chatID)
	} else {
		addReaction(c.bot, callbackQuery.Message, "👎")
	}

	chatID = tu.ID(annInfo.ChatID)
	msg := fmt.Sprintf("Объявление с номером %d удалено администратором. Смотрите правила /start", annInfo.Id)
	message := tu.Message(chatID, msg)
	c.bot.SendMessage(message)

}
