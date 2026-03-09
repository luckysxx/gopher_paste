package repository

import (
	"context"
	"database/sql"
	"project/common/dberr"
	"project/services/paste/db"
	"project/services/paste/model"
	"time"
)

type PasteRepository interface {
	Create(ctx context.Context, ownerID int64, params *model.CreatePasteRequest, shortLink string) (*model.PasteResponse, error)
	GetByID(ctx context.Context, id int64) (*model.PasteResponse, error)
	ListByOwner(ctx context.Context, ownerID int64) ([]model.PasteResponse, error)
	Update(ctx context.Context, ownerID, id int64, params *model.UpdatePasteRequest) (*model.PasteResponse, error)
}

type pasteRepository struct {
	q *db.Queries
}

func NewPasteRepository(conn *sql.DB) PasteRepository {
	return &pasteRepository{q: db.New(conn)}
}

func (r *pasteRepository) Create(ctx context.Context, ownerID int64, params *model.CreatePasteRequest, shortLink string) (*model.PasteResponse, error) {
	paste, err := r.q.CreatePaste(ctx, db.CreatePasteParams{
		OwnerID: ownerID,
		Title:   params.Title,
		ShortLink: sql.NullString{
			String: shortLink,
			Valid:  shortLink != "",
		},
		Content:    params.Content,
		Language:   params.Language,
		Visibility: resolveVisibility(params.Visibility),
	})
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toResponse(&paste), nil
}

func (r *pasteRepository) GetByID(ctx context.Context, id int64) (*model.PasteResponse, error) {
	paste, err := r.q.GetPaste(ctx, id)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toResponse(&paste), nil
}

func (r *pasteRepository) ListByOwner(ctx context.Context, ownerID int64) ([]model.PasteResponse, error) {
	rows, err := r.q.ListMyPastes(ctx, ownerID)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	results := make([]model.PasteResponse, 0)
	for i := range rows {
		resp := toResponse(&rows[i])
		results = append(results, *resp)
	}

	return results, nil
}

func (r *pasteRepository) Update(ctx context.Context, ownerID, id int64, params *model.UpdatePasteRequest) (*model.PasteResponse, error) {
	paste, err := r.q.UpdatePaste(ctx, db.UpdatePasteParams{
		ID:         id,
		Title:      params.Title,
		Content:    params.Content,
		Language:   params.Language,
		Visibility: resolveVisibility(params.Visibility),
		OwnerID:    ownerID,
	})
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}

	return toResponse(&paste), nil
}

func resolveVisibility(visibility string) string {
	if visibility == "public" {
		return "public"
	}
	return "private"
}

func toResponse(paste *db.Paste) *model.PasteResponse {
	resp := &model.PasteResponse{
		ID:         paste.ID,
		OwnerID:    paste.OwnerID,
		Title:      paste.Title,
		Content:    paste.Content,
		Language:   paste.Language,
		Visibility: paste.Visibility,
		CreatedAt:  paste.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  paste.UpdatedAt.Format(time.RFC3339),
	}

	if paste.ShortLink.Valid {
		resp.ShortLink = paste.ShortLink.String
	}

	return resp
}
