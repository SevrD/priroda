package users

import (
	"context"
	"database/sql"
	"main/internal/domain"
	"main/internal/queries"

	"time"

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

func (c *Users) Register(ctx context.Context, tgID int64, login string, name string, date int, chatID int64) error {

	var params queries.CreateUserParams

	sqlTgID := sql.NullInt64{Int64: tgID, Valid: true}

	sqlTimeStamp := sql.NullTime{Time: time.Unix(int64(date), 0), Valid: true}

	sqlChatID := sql.NullInt64{Int64: chatID, Valid: true}

	params.Tgid = sqlTgID
	params.Login = login
	params.Name = name
	params.Createdata = sqlTimeStamp
	params.Chatid = sqlChatID

	return c.queries.CreateUser(ctx, params)

}

func (c *Users) GetUserInfo(ctx context.Context, tgID int64) (name string, login string, ban bool, err error) {

	sqlTgID := sql.NullInt64{Int64: tgID, Valid: true}

	userInfo, err := c.queries.GetUserInfo(ctx, sqlTgID)
	return userInfo.Name, userInfo.Login, userInfo.Ban.Bool, err
}

func (c *Users) Ban(ctx context.Context, tgID int64) error {

	sqlTgID := sql.NullInt64{Int64: tgID, Valid: true}

	return c.queries.Ban(ctx, sqlTgID)

}

func (c *Users) UnBan(ctx context.Context, tgID int64) error {

	sqlTgID := sql.NullInt64{Int64: tgID, Valid: true}
	return c.queries.UnBan(ctx, sqlTgID)

}

func (c *Users) SendError(chatID telego.ChatID) {
	errorText := "Что-то пошло не так. Попробуйте повторить операцию или обратитесь в техническую поддержку @" + c.loginAdmin
	message := tu.Message(chatID, errorText)
	c.bot.SendMessage(message)
}
