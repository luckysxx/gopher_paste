package repository

import (
	"context"
	"project/common/dberr"
	"project/services/paste/db"
)

type PasteRepository interface {
	Create(ctx context.Context, params *db.CreatePasteParams) (*db.Paste, error)
	GetByShortLink(ctx context.Context, shortLink string) (*db.Paste, error)
}

type pasteRepository struct {
	q *db.Queries
}

func NewPasteRepository(q *db.Queries) PasteRepository {
	return &pasteRepository{q: q}
}

func (r *pasteRepository) Create(ctx context.Context, params *db.CreatePasteParams) (*db.Paste, error) {
	paste, err := r.q.CreatePaste(ctx, *params)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}
	return &paste, nil
}

func (r *pasteRepository) GetByShortLink(ctx context.Context, shortLink string) (*db.Paste, error) {
	paste, err := r.q.GetPaste(ctx, shortLink)
	if err != nil {
		return nil, dberr.ParseDBError(err)
	}
	return &paste, nil
}
