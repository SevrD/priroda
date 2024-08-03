// Package domain for interfaces
package domain

import (
	"context"
	"main/internal/models"

	"github.com/mymmrac/telego"
)

type Users interface {
	Register(ctx context.Context, tgid int64, login string, name string, date int, chatID int64) error
	GetUserInfo(ctx context.Context, tgID int64) (name string, login string, ban bool, err error)
	Ban(ctx context.Context, tgID int64) error
	UnBan(ctx context.Context, tgID int64) error
	SendError(hatID telego.ChatID)
	LoginAdmin() string
}

type ChatStatus interface {
	Save(ctx context.Context, tgid int64, statusCode models.StatusCode, ann_id int64) error
	Get(ctx context.Context, tgid int64) models.StatusCode
	GetAnnId(ctx context.Context, tgid int64) int64
}

type Announcement interface {
	Add(ctx context.Context, tgid int64, txt string, chatID int64) (int64, error)
	GetAnnouncement(ctx context.Context, tgid int64, id int64) (txt string, publicID int64, err error)
	SetAdminMsgID(ctx context.Context, id int64, adm_msg_id int64) error
	SetPublicID(ctx context.Context, id int64, publicID int64) error
	AddPhoto(ctx context.Context, annID int64, fileID string) error
	GetAnnouncementOnAdmMsgID(ctx context.Context, admMsgID int64) (annInfo *models.AnnouncementInfo, err error)
}

type CallBacker interface {
	CallBack(ctx context.Context, callbackQuery *telego.CallbackQuery)
}

type Core interface {
	SendMessageWithButtons(chatID telego.ChatID, text string, buttons []models.Button)
	SendMessageAfterAddAnnouncement(ctx context.Context, tgUserID int64, chatID telego.ChatID, fileID *string)
	SendDefaultMessage(chatID telego.ChatID)
	SaveAdminChatID(chatID int64)
	AddContacts(ctx context.Context, tgUserID int64, txt string) (string, error)
	Contacts(ctx context.Context, tgUserID int64) string
	SendDeleteRequest(ctx context.Context, tgUserID int64, annID int64, chatID telego.ChatID) error
}

type CommandsWorker interface {
	Start(ctx context.Context, chatID telego.ChatID, tgID int64, rules string)
	Add(ctx context.Context, chatID telego.ChatID, tgID int64)
	Delete(ctx context.Context, chatID telego.ChatID, tgID int64)
}

type Messager interface {
	Text(ctx context.Context, message *telego.Message)
	Photo(ctx context.Context, message *telego.Message)
	Number(ctx context.Context, message *telego.Message)
}
