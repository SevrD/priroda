package users

import (
	"context"
	"main/internal/domain"
	"main/internal/queries"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Users struct {
	queries    *queries.Queries
	loginAdmin string
	bot        *telego.Bot
}

func NewUsersClient(queries *queries.Queries, loginAdmin string, bot *telego.Bot) domain.Users {
	return &Users{
		queries:    queries,
		loginAdmin: loginAdmin,
		bot:        bot,
	}
}

func (c *Users) Register(ctx context.Context, tgid int64, login string, name string, date int, chatID int64) error {

	var params queries.CreateUserParams

	var timeStamp pgtype.Timestamp
	timeStamp.Scan(time.Unix(int64(date), 0))

	var pgChatID pgtype.Int8
	pgChatID.Scan(chatID)

	var pgTgid pgtype.Int8
	pgTgid.Scan(tgid)

	params.Tgid = pgTgid
	params.Login = login
	params.Name = name
	params.Createdata = timeStamp
	params.Chatid = pgChatID

	_, err := c.queries.CreateUser(ctx, params)

	return err
}

func (c *Users) GetUserInfo(ctx context.Context, tgID int64) (name string, login string, ban bool, err error) {

	var pgTgID pgtype.Int8
	pgTgID.Scan(tgID)

	userInfo, err := c.queries.GetUserInfo(ctx, pgTgID)
	return userInfo.Name, userInfo.Login, userInfo.Ban.Bool, err
}

func (c *Users) Ban(ctx context.Context, tgID int64) error {

	var pgTgID pgtype.Int8
	pgTgID.Scan(tgID)

	return c.queries.Ban(ctx, pgTgID)

}

func (c *Users) UnBan(ctx context.Context, tgID int64) error {

	var pgTgID pgtype.Int8
	pgTgID.Scan(tgID)

	return c.queries.UnBan(ctx, pgTgID)

}

func (c *Users) SendError(chatID telego.ChatID) {
	errorText := "Что-то пошло не так. Попробуйте повторить операцию или обратитесь в техническую поддержку @" + c.loginAdmin
	message := tu.Message(chatID, errorText)
	c.bot.SendMessage(message)
}
