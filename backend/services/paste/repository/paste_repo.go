package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"project/services/paste/db"
)

var ErrPasteNotFound = errors.New("paste not found")

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
		return nil, fmt.Errorf("创建 paste 失败: %w", err)
	}
	return &paste, nil
}

func (r *pasteRepository) GetByShortLink(ctx context.Context, shortLink string) (*db.Paste, error) {
	paste, err := r.q.GetPaste(ctx, shortLink)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPasteNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询 paste 失败 (shortLink=%s): %w", shortLink, err)
	}
	return &paste, nil
}
