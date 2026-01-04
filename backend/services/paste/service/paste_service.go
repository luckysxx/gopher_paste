package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"project/common/cache"
	"project/services/paste/db"
	"project/services/paste/model"
	"project/services/paste/repository"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

// 领域错误定义（协议无关）
var (
	ErrShortLinkGeneration = errors.New("无法生成唯一 shortlink")
	ErrPasteNotFound       = repository.ErrPasteNotFound // 复用 Repository 的错误
)

type PasteService interface {
	Create(ctx context.Context, req *model.CreatePasteRequest) (*db.Paste, error)
	GetByShortLink(ctx context.Context, shortLink string) (*db.Paste, error)
}

type pasteService struct {
	repo   repository.PasteRepository
	cache  cache.Cache
	logger *zap.Logger
}

func NewPasteService(repo repository.PasteRepository, cache cache.Cache, logger *zap.Logger) PasteService {
	return &pasteService{repo: repo, cache: cache, logger: logger}
}

func (s *pasteService) Create(ctx context.Context, req *model.CreatePasteRequest) (*db.Paste, error) {
	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		cand, err := generateShortLink(8)
		if err != nil {
			s.logger.Error("生成 shortlink 失败", zap.Error(err))
			return nil, err
		}

		params := &db.CreatePasteParams{
			ShortLink: cand,
			Content:   req.Content,
			Language:  req.Language,
			ExpiresAt: sql.NullTime{Valid: false},
		}

		paste, err := s.repo.Create(ctx, params)
		if err == nil {
			return paste, nil
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			// 唯一键冲突，重试
			s.logger.Warn("shortlink 冲突，重试生成", zap.String("shortLink", cand))
			continue
		}

		// 其他数据库错误，返回标准 error
		s.logger.Error("创建 paste 失败",
			zap.String("shortLink", cand),
			zap.Error(err),
		)
		return nil, fmt.Errorf("数据库错误: %w", err)
	}

	s.logger.Error("无法生成唯一 shortlink")
	return nil, ErrShortLinkGeneration
}

func generateShortLink(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		b[i] = letters[idx.Int64()]
	}
	return string(b), nil
}

func (s *pasteService) GetByShortLink(ctx context.Context, shortLink string) (*db.Paste, error) {
	cacheKey := "paste:" + shortLink

	// 1. 查缓存
	if cached, err := s.getFromCache(ctx, cacheKey); err == nil {
		s.logger.Info("Cache hit", zap.String("key", shortLink))
		return cached, nil
	}

	// 2. 查数据库
	paste, err := s.repo.GetByShortLink(ctx, shortLink)
	if err != nil {
		// Repository 返回的是领域错误
		if errors.Is(err, repository.ErrPasteNotFound) {
			// 直接返回，由 Handler 层处理
			return nil, ErrPasteNotFound
		}
		// 其他数据库错误
		s.logger.Error("查询 paste 失败",
			zap.String("shortLink", shortLink),
			zap.Error(err),
		)
		return nil, fmt.Errorf("数据库错误: %w", err)
	}

	// 3. 回写缓存（异步或忽略错误）
	s.setCache(ctx, cacheKey, paste)
	s.logger.Info("Cache miss - loaded from DB", zap.String("key", shortLink))

	return paste, nil
}

func (s *pasteService) getFromCache(ctx context.Context, key string) (*db.Paste, error) {
	val, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var paste db.Paste
	if err := json.Unmarshal([]byte(val), &paste); err != nil {
		return nil, err
	}
	return &paste, nil
}

func (s *pasteService) setCache(ctx context.Context, key string, paste *db.Paste) {
	data, _ := json.Marshal(paste)
	if err := s.cache.Set(ctx, key, string(data), time.Hour); err != nil {
		s.logger.Warn("Failed to set cache", zap.Error(err))
	}
}
