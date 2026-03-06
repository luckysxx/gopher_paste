package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"project/common/dberr"
	"project/services/paste/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
)

func newMockPasteRepository(t *testing.T) (PasteRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建 sqlmock 失败: %v", err)
	}

	repo := NewPasteRepository(db.New(sqlDB))
	cleanup := func() {
		sqlDB.Close()
	}

	return repo, mock, cleanup
}

func TestPasteRepository_Create_Success(t *testing.T) {
	repo, mock, cleanup := newMockPasteRepository(t)
	defer cleanup()

	createdAt := time.Now()
	expiresAt := sql.NullTime{Valid: false}

	params := &db.CreatePasteParams{
		ShortLink: "abc12345",
		Content:   "hello",
		Language:  "go",
		ExpiresAt: expiresAt,
	}

	mock.ExpectQuery("INSERT INTO pastes").
		WithArgs(params.ShortLink, params.Content, params.Language, params.ExpiresAt).
		WillReturnRows(sqlmock.NewRows([]string{"short_link", "content", "language", "expires_at", "created_at"}).
			AddRow(params.ShortLink, params.Content, params.Language, nil, createdAt))

	paste, err := repo.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("Create() 返回错误: %v", err)
	}

	if paste.ShortLink != params.ShortLink || paste.Content != params.Content || paste.Language != params.Language {
		t.Fatalf("Create() 返回值不符合预期: %+v", paste)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 期望未满足: %v", err)
	}
}

func TestPasteRepository_Create_MapsDBError(t *testing.T) {
	repo, mock, cleanup := newMockPasteRepository(t)
	defer cleanup()

	params := &db.CreatePasteParams{
		ShortLink: "dup12345",
		Content:   "hello",
		Language:  "go",
		ExpiresAt: sql.NullTime{Valid: false},
	}

	mock.ExpectQuery("INSERT INTO pastes").
		WithArgs(params.ShortLink, params.Content, params.Language, params.ExpiresAt).
		WillReturnError(&pq.Error{Code: dberr.PgErrUniqueViolation, Constraint: "pastes_short_link_key"})

	paste, err := repo.Create(context.Background(), params)
	if paste != nil {
		t.Fatal("Create() 出错时应返回 nil paste")
	}
	if !errors.Is(err, dberr.ErrShortLinkDuplicate) {
		t.Fatalf("期望 ErrShortLinkDuplicate, 实际: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 期望未满足: %v", err)
	}
}

func TestPasteRepository_GetByShortLink_Success(t *testing.T) {
	repo, mock, cleanup := newMockPasteRepository(t)
	defer cleanup()

	createdAt := time.Now()
	expiresAt := time.Now().Add(2 * time.Hour)
	shortLink := "k9Zx1Qwe"

	mock.ExpectQuery("SELECT short_link, content, language, expires_at, created_at FROM pastes").
		WithArgs(shortLink).
		WillReturnRows(sqlmock.NewRows([]string{"short_link", "content", "language", "expires_at", "created_at"}).
			AddRow(shortLink, "world", "rust", expiresAt, createdAt))

	paste, err := repo.GetByShortLink(context.Background(), shortLink)
	if err != nil {
		t.Fatalf("GetByShortLink() 返回错误: %v", err)
	}

	if paste.ShortLink != shortLink || paste.Content != "world" || paste.Language != "rust" {
		t.Fatalf("GetByShortLink() 返回值不符合预期: %+v", paste)
	}

	if !paste.ExpiresAt.Valid {
		t.Fatal("ExpiresAt 应该是有效时间")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 期望未满足: %v", err)
	}
}

func TestPasteRepository_GetByShortLink_NotFound(t *testing.T) {
	repo, mock, cleanup := newMockPasteRepository(t)
	defer cleanup()

	shortLink := "missing"
	mock.ExpectQuery("SELECT short_link, content, language, expires_at, created_at FROM pastes").
		WithArgs(shortLink).
		WillReturnError(sql.ErrNoRows)

	paste, err := repo.GetByShortLink(context.Background(), shortLink)
	if paste != nil {
		t.Fatal("GetByShortLink() 未找到时应返回 nil paste")
	}
	if !errors.Is(err, dberr.ErrNoRows) {
		t.Fatalf("期望 ErrNoRows, 实际: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("SQL 期望未满足: %v", err)
	}
}
