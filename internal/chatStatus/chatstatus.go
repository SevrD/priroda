package chatstatus

import (
	"context"
	"database/sql"
	"log"
	"main/internal/domain"
	"main/internal/models"
	"main/internal/queries"
)

type ChatStatus struct {
	queries *queries.Queries
}

func NewChatStatus(queries *queries.Queries) domain.ChatStatus {
	return &ChatStatus{
		queries: queries,
	}
}

func (c *ChatStatus) Save(ctx context.Context, tgID int64, statusCode models.StatusCode, annID int64) error {

	var params queries.SetStatusParams

	sqlTgID := sql.NullInt64{Int64: tgID, Valid: true}

	sqlPgStatus := sql.NullInt64{Int64: int64(statusCode), Valid: true}

	sqlAnnID := sql.NullInt64{Int64: annID, Valid: true}

	params.Tgid = sqlTgID
	params.Status = sqlPgStatus
	params.Annid = sqlAnnID

	return c.queries.SetStatus(ctx, params)

}

func (c *ChatStatus) Get(ctx context.Context, tgID int64) models.StatusCode {

	sqlTgID := sql.NullInt64{Int64: tgID, Valid: true}

	status, err := c.queries.GetStatus(ctx, sqlTgID)

	if err != nil && err != sql.ErrNoRows {
		log.Println("Get status error:", err)
	}

	return models.StatusCode(status.Int64)

}

func (c *ChatStatus) GetAnnId(ctx context.Context, tgID int64) int64 {

	sqlTgID := sql.NullInt64{Int64: tgID, Valid: true}

	id, err := c.queries.GetAnnId(ctx, sqlTgID)

	if err != nil && err != sql.ErrNoRows {
		log.Println("Get announcement id:", err)
	}

	return id.Int64

}
