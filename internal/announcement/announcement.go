package announcement

import (
	"context"
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

	var pgTgid pgtype.Int8
	pgTgid.Scan(tgid)

	params.Tgid = pgTgid

	var pgTxt pgtype.Text
	pgTxt.Scan(txt)

	params.Txt = pgTxt

	var pgChatID pgtype.Int8
	pgChatID.Scan(chatID)

	params.Chatid = pgChatID

	return c.queries.AddAnnouncement(ctx, params)

}

func (c *Announcement) GetAnnouncement(ctx context.Context, tgid int64, ann_id int64) (txt string, publicID int64, err error) {

	var pgTgid pgtype.Int8
	pgTgid.Scan(tgid)

	GetAnnouncementParams := queries.GetAnnouncementParams{
		Tgid: pgTgid,
		ID:   ann_id,
	}

	result, err := c.queries.GetAnnouncement(ctx, GetAnnouncementParams)

	if err != nil {
		return "", 0, err
	}

	return result.Txt.String, result.Publicid.Int64, nil

}

func (c *Announcement) SetAdminMsgID(ctx context.Context, id int64, adm_msg_id int64) error {

	var pgAdmMsgID pgtype.Int8
	pgAdmMsgID.Scan(adm_msg_id)

	params := queries.SetAdminMsgIDParams{
		ID:       id,
		Admmsgid: pgAdmMsgID,
	}

	return c.queries.SetAdminMsgID(ctx, params)
}

func (c *Announcement) AddPhoto(ctx context.Context, annID int64, fileID string) error {

	var pgTxt pgtype.Text
	pgTxt.Scan(fileID)

	photoParams := queries.AddPhotoParams{ID: annID, Fileid: pgTxt}

	return c.queries.AddPhoto(ctx, photoParams)
}

func (c *Announcement) GetAnnouncementOnAdmMsgID(ctx context.Context, admMsgID int64) (*models.AnnouncementInfo, error) {

	var pgAdmMsgID pgtype.Int8
	pgAdmMsgID.Scan(admMsgID)

	result, err := c.queries.GetAnnouncementOnAdmMsgID(ctx, pgAdmMsgID)

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

	var pgPublicID pgtype.Int8
	pgPublicID.Scan(publicID)

	params := queries.SetPublicIDParams{
		ID:       id,
		Publicid: pgPublicID,
	}

	return c.queries.SetPublicID(ctx, params)

}
