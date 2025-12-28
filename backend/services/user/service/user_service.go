package service

import (
	"context"
	"fmt"
	"project/common/auth"
	"project/services/user/db"
	"project/services/user/model"
	"project/services/user/repository"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{repo: repo, logger: logger}
}

func (s *userService) Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
	// 加密密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("密码加密失败", zap.Error(err))
		return nil, err
	}

	// 构造参数
	params := &db.CreateUserParams{
		Username: req.Username,
		Password: string(hashedPwd),
		Email:    req.Email,
	}

	// 调用数据库创建用户
	user, err := s.repo.Create(ctx, params)
	if err != nil {
		return nil, err
	}

	resp := &model.RegisterResponse{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return resp, nil
}

func (s *userService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	// 根据用户名获取用户
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 比较密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	// 生成 JWT Token
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("生成 Token 失败: %v", err)
	}

	// 生成登录响应
	resp := &model.LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return resp, nil
}
