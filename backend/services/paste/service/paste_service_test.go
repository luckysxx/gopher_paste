package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"project/services/paste/db"

	"go.uber.org/zap"
)

type mockPasteRepo struct {
	createFn   func(ctx context.Context, params *db.CreatePasteParams) (*db.Paste, error)
	getFn      func(ctx context.Context, shortLink string) (*db.Paste, error)
	createCall int
	getCall    int
}

func (m *mockPasteRepo) Create(ctx context.Context, params *db.CreatePasteParams) (*db.Paste, error) {
	m.createCall++
	if m.createFn != nil {
		return m.createFn(ctx, params)
	}
	return nil, errors.New("createFn not implemented")
}

func (m *mockPasteRepo) GetByShortLink(ctx context.Context, shortLink string) (*db.Paste, error) {
	m.getCall++
	if m.getFn != nil {
		return m.getFn(ctx, shortLink)
	}
	return nil, errors.New("getFn not implemented")
}

type mockCache struct {
	getFn      func(ctx context.Context, key string) (string, error)
	setFn      func(ctx context.Context, key string, value string, exp time.Duration) error
	delFn      func(ctx context.Context, keys ...string) error
	setCall    int
	lastSetKey string
	lastSetVal string
}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	if m.getFn != nil {
		return m.getFn(ctx, key)
	}
	return "", errors.New("getFn not implemented")
}

func (m *mockCache) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	m.setCall++
	m.lastSetKey = key
	m.lastSetVal = value
	if m.setFn != nil {
		return m.setFn(ctx, key, value, exp)
	}
	return nil
}

func (m *mockCache) Del(ctx context.Context, keys ...string) error {
	if m.delFn != nil {
		return m.delFn(ctx, keys...)
	}
	return nil
}

func TestNewPasteService(t *testing.T) {
	svc := NewPasteService(&mockPasteRepo{}, &mockCache{}, zap.NewNop())
	if svc == nil {
		t.Fatal("NewPasteService() 返回了 nil")
	}
}

func TestGetByShortLink_CacheHit(t *testing.T) {
	expected := &db.Paste{
		ShortLink: "abc12345",
		Content:   "hello",
		Language:  "go",
	}
	buf, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("准备缓存数据失败: %v", err)
	}

	repo := &mockPasteRepo{
		getFn: func(ctx context.Context, shortLink string) (*db.Paste, error) {
			t.Fatal("缓存命中时不应访问 repo")
			return nil, nil
		},
	}
	cache := &mockCache{
		getFn: func(ctx context.Context, key string) (string, error) {
			if key != "paste:abc12345" {
				t.Fatalf("缓存 key 不对, got=%s", key)
			}
			return string(buf), nil
		},
	}

	svc := NewPasteService(repo, cache, zap.NewNop())
	got, err := svc.GetByShortLink(context.Background(), "abc12345")
	if err != nil {
		t.Fatalf("GetByShortLink() 返回错误: %v", err)
	}

	if got.ShortLink != expected.ShortLink || got.Content != expected.Content || got.Language != expected.Language {
		t.Fatalf("缓存命中返回值不符合预期, got=%+v", got)
	}
	if repo.getCall != 0 {
		t.Fatalf("缓存命中不应访问 repo, 实际访问次数=%d", repo.getCall)
	}
}

func TestGetByShortLink_CacheMissFallsBackToRepoAndSetCache(t *testing.T) {
	expected := &db.Paste{
		ShortLink: "k9Zx1Qwe",
		Content:   "world",
		Language:  "rust",
		ExpiresAt: sql.NullTime{Valid: false},
	}

	repo := &mockPasteRepo{
		getFn: func(ctx context.Context, shortLink string) (*db.Paste, error) {
			if shortLink != "k9Zx1Qwe" {
				t.Fatalf("repo 查询参数不对, got=%s", shortLink)
			}
			return expected, nil
		},
	}
	cache := &mockCache{
		getFn: func(ctx context.Context, key string) (string, error) {
			return "", errors.New("cache miss")
		},
	}

	svc := NewPasteService(repo, cache, zap.NewNop())
	got, err := svc.GetByShortLink(context.Background(), "k9Zx1Qwe")
	if err != nil {
		t.Fatalf("GetByShortLink() 返回错误: %v", err)
	}

	if got.ShortLink != expected.ShortLink {
		t.Fatalf("返回 shortLink 不对, 期望=%s 实际=%s", expected.ShortLink, got.ShortLink)
	}
	if repo.getCall != 1 {
		t.Fatalf("缓存未命中应访问 repo 1 次, 实际=%d", repo.getCall)
	}
	if cache.setCall != 1 {
		t.Fatalf("缓存未命中后应回写缓存 1 次, 实际=%d", cache.setCall)
	}
	if cache.lastSetKey != "paste:k9Zx1Qwe" {
		t.Fatalf("回写缓存 key 不对, got=%s", cache.lastSetKey)
	}

}
