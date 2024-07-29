package announcement

import (
	"context"
	"database/sql"
	"main/internal/domain"
	"main/internal/models"
	"main/internal/queries"

	"github.com/jackc/pgx/v5/pgtype"
)

type Announcement struct {
	queries *queries.Queries
}

func NewAnnouncement(queries *queries.Queries) domain.Announcement {
	return &Announcement{
		queries: queries,
	}
}

func (c *Announcement) Add(ctx context.Context, tgid int64, txt string, chatID int64) (int64, error) {

	var params queries.AddAnnouncementParams

	var sqlID sql.NullInt64
	sqlID.Scan(tgid)

	params.Tgid = sqlID

	var sqlTxt sql.NullString
	sqlTxt.Scan(txt)

	params.Txt = sqlTxt

	var sqlChatID sql.NullInt64
	sqlChatID.Scan(chatID)

	params.Chatid = sqlChatID

	return c.queries.AddAnnouncement(ctx, params)

}

func (c *Announcement) GetAnnouncement(ctx context.Context, tgID int64, annID int64) (txt string, publicID int64, err error) {

	var pgTgid pgtype.Int8
	pgTgid.Scan(tgID)

	sqlID := sql.NullInt64{Int64: tgID, Valid: true}

	GetAnnouncementParams := queries.GetAnnouncementParams{
		Tgid: sqlID,
		ID:   annID,
	}

	result, err := c.queries.GetAnnouncement(ctx, GetAnnouncementParams)

	if err != nil {
		return "", 0, err
	}

	return result.Txt.String, result.Publicid.Int64, nil

}

func (c *Announcement) SetAdminMsgID(ctx context.Context, id int64, admMsgID int64) error {

	sqlAdmMsgID := sql.NullInt64{Int64: admMsgID, Valid: true}

	params := queries.SetAdminMsgIDParams{
		ID:       id,
		Admmsgid: sqlAdmMsgID,
	}

	return c.queries.SetAdminMsgID(ctx, params)
}

func (c *Announcement) AddPhoto(ctx context.Context, annID int64, fileID string) error {

	sqlAdmMsgID := sql.NullString{String: fileID, Valid: true}

	photoParams := queries.AddPhotoParams{ID: annID, Fileid: sqlAdmMsgID}

	return c.queries.AddPhoto(ctx, photoParams)
}

func (c *Announcement) GetAnnouncementOnAdmMsgID(ctx context.Context, admMsgID int64) (*models.AnnouncementInfo, error) {

	sqlAdmMsgID := sql.NullInt64{Int64: admMsgID, Valid: true}

	result, err := c.queries.GetAnnouncementOnAdmMsgID(ctx, sqlAdmMsgID)

	if err != nil {
		return nil, err
	}

	annInfo := &models.AnnouncementInfo{
		Text:     result.Txt.String,
		FileID:   result.Fileid.String,
		ChatID:   result.Chatid.Int64,
		Id:       result.ID,
		TgID:     result.Tgid.Int64,
		PublicID: result.Publicid.Int64,
	}

	return annInfo, nil
}

func (c *Announcement) SetPublicID(ctx context.Context, id int64, publicID int64) error {

	sqlPublicID := sql.NullInt64{Int64: publicID, Valid: true}

	params := queries.SetPublicIDParams{
		ID:       id,
		Publicid: sqlPublicID,
	}

	return c.queries.SetPublicID(ctx, params)

}
