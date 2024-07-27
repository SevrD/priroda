package chatstatus

import (
	"context"
	"database/sql"
	"log"
	"main/internal/domain"
	"main/internal/models"
	"main/internal/queries"

	"github.com/jackc/pgx/v5/pgtype"
)

type ChatStatus struct {
	queries *queries.Queries
}

func NewChatStatus(queries *queries.Queries) domain.ChatStatus {
	return &ChatStatus{
		queries: queries,
	}
}

func (c *ChatStatus) Save(ctx context.Context, tgid int64, statusCode models.StatusCode, ann_id int64) error {

	var params queries.SetStatusParams
	var pgTgid pgtype.Int8
	err := pgTgid.Scan(tgid)

	if err != nil {
		return err
	}

	var pgStatus pgtype.Int8
	err = pgStatus.Scan(int64(statusCode))

	if err != nil {
		return err
	}

	var pgAnnid pgtype.Int8
	err = pgAnnid.Scan(ann_id)

	if err != nil {
		return err
	}

	params.Tgid = pgTgid
	params.Status = pgStatus
	params.Annid = pgAnnid

	_, err = c.queries.SetStatus(ctx, params)
	if err != nil {
		log.Println("Save status error:", err)
	}

	return err
}

func (c *ChatStatus) Get(ctx context.Context, tgid int64) models.StatusCode {

	var pgTgid pgtype.Int8
	err := pgTgid.Scan(tgid)

	if err != nil {
		log.Println("Value error:", err)
	}

	status, err := c.queries.GetStatus(ctx, pgTgid)

	if err != nil && err != sql.ErrNoRows {
		log.Println("Get status error:", err)
	}

	return models.StatusCode(status.Int64)

}

func (c *ChatStatus) GetAnnId(ctx context.Context, tgid int64) int64 {

	var pgTgid pgtype.Int8
	err := pgTgid.Scan(tgid)

	if err != nil {
		log.Println("Value error:", err)
	}

	id, err := c.queries.GetAnnId(ctx, pgTgid)

	if err != nil && err != sql.ErrNoRows {
		log.Println("Get announcement id:", err)
	}

	return id.Int64

}
