package service

import (
	"context"
	"errors"
	"testing"

	"project/common/dberr"
	"project/services/paste/model"

	"go.uber.org/zap"
)

type mockPasteRepo struct {
	createFn      func(ctx context.Context, ownerID int64, params *model.CreatePasteRequest, shortLink string) (*model.PasteResponse, error)
	getByIDFn     func(ctx context.Context, id int64) (*model.PasteResponse, error)
	listByOwnerFn func(ctx context.Context, ownerID int64) ([]model.PasteResponse, error)
	updateFn      func(ctx context.Context, ownerID, id int64, params *model.UpdatePasteRequest) (*model.PasteResponse, error)
}

func (m *mockPasteRepo) Create(ctx context.Context, ownerID int64, params *model.CreatePasteRequest, shortLink string) (*model.PasteResponse, error) {
	if m.createFn != nil {
		return m.createFn(ctx, ownerID, params, shortLink)
	}
	return nil, errors.New("createFn not implemented")
}

func (m *mockPasteRepo) GetByID(ctx context.Context, id int64) (*model.PasteResponse, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("getByIDFn not implemented")
}

func (m *mockPasteRepo) ListByOwner(ctx context.Context, ownerID int64) ([]model.PasteResponse, error) {
	if m.listByOwnerFn != nil {
		return m.listByOwnerFn(ctx, ownerID)
	}
	return nil, errors.New("listByOwnerFn not implemented")
}

func (m *mockPasteRepo) Update(ctx context.Context, ownerID, id int64, params *model.UpdatePasteRequest) (*model.PasteResponse, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, ownerID, id, params)
	}
	return nil, errors.New("updateFn not implemented")
}

func TestNewPasteService(t *testing.T) {
	svc := NewPasteService(&mockPasteRepo{}, zap.NewNop())
	if svc == nil {
		t.Fatal("NewPasteService() 返回了 nil")
	}
}

func TestListMine(t *testing.T) {
	repo := &mockPasteRepo{
		listByOwnerFn: func(ctx context.Context, ownerID int64) ([]model.PasteResponse, error) {
			if ownerID != 7 {
				t.Fatalf("ownerID 不匹配: %d", ownerID)
			}
			return []model.PasteResponse{{ID: 1, OwnerID: 7, Title: "hello"}}, nil
		},
	}

	svc := NewPasteService(repo, zap.NewNop())
	list, err := svc.ListMine(context.Background(), 7)
	if err != nil {
		t.Fatalf("ListMine() 返回错误: %v", err)
	}

	if len(list) != 1 || list[0].Title != "hello" {
		t.Fatalf("ListMine() 返回值不符合预期: %+v", list)
	}
}

func TestUpdate_ForbiddenWhenOwnerMismatch(t *testing.T) {
	repo := &mockPasteRepo{
		getByIDFn: func(ctx context.Context, id int64) (*model.PasteResponse, error) {
			return &model.PasteResponse{ID: id, OwnerID: 100, Title: "old"}, nil
		},
	}

	svc := NewPasteService(repo, zap.NewNop())
	_, err := svc.Update(context.Background(), 101, 1, &model.UpdatePasteRequest{
		Title:    "new",
		Content:  "code",
		Language: "go",
	})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("期望 ErrForbidden, 实际: %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	repo := &mockPasteRepo{
		getByIDFn: func(ctx context.Context, id int64) (*model.PasteResponse, error) {
			return &model.PasteResponse{ID: id, OwnerID: 7, Title: "old"}, nil
		},
		updateFn: func(ctx context.Context, ownerID, id int64, params *model.UpdatePasteRequest) (*model.PasteResponse, error) {
			if ownerID != 7 || id != 3 {
				t.Fatalf("更新参数不匹配 ownerID=%d id=%d", ownerID, id)
			}
			return &model.PasteResponse{ID: id, OwnerID: ownerID, Title: params.Title}, nil
		},
	}

	svc := NewPasteService(repo, zap.NewNop())
	res, err := svc.Update(context.Background(), 7, 3, &model.UpdatePasteRequest{
		Title:    "new-title",
		Content:  "code",
		Language: "go",
	})
	if err != nil {
		t.Fatalf("Update() 返回错误: %v", err)
	}
	if res.Title != "new-title" {
		t.Fatalf("Update() 返回值不符合预期: %+v", res)
	}
}

func TestUpdate_RepoNotFoundMapsForbidden(t *testing.T) {
	repo := &mockPasteRepo{
		getByIDFn: func(ctx context.Context, id int64) (*model.PasteResponse, error) {
			return &model.PasteResponse{ID: id, OwnerID: 9, Title: "old"}, nil
		},
		updateFn: func(ctx context.Context, ownerID, id int64, params *model.UpdatePasteRequest) (*model.PasteResponse, error) {
			return nil, dberr.ErrNoRows
		},
	}

	svc := NewPasteService(repo, zap.NewNop())
	_, err := svc.Update(context.Background(), 9, 1, &model.UpdatePasteRequest{Title: "n", Content: "c", Language: "go"})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("期望 ErrForbidden, 实际: %v", err)
	}
}
