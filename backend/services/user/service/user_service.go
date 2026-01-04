package service

import (
	"context"
	"errors"
	"fmt"
	"project/common/auth"
	"project/services/user/db"
	"project/services/user/model"
	"project/services/user/repository"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 领域错误定义（协议无关）
var (
	ErrUserNotFound       = repository.ErrUserNotFound // 复用 Repository 的错误
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrTokenGeneration    = errors.New("生成 Token 失败")
)

type UserService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}

type userService struct {
	repo       repository.UserRepository
	logger     *zap.Logger
	jwtManager *auth.JWTManager
}

func NewUserService(repo repository.UserRepository, jwtManager *auth.JWTManager, logger *zap.Logger) UserService {
	return &userService{repo: repo, jwtManager: jwtManager, logger: logger}
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
		s.logger.Error("创建用户失败", zap.Error(err))
		return nil, fmt.Errorf("数据库错误: %w", err)
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
		if errors.Is(err, repository.ErrUserNotFound) {
			// 返回领域错误，不泄露具体是用户名错误
			return nil, ErrInvalidCredentials
		}
		// 其他数据库错误
		s.logger.Error("查询用户失败", zap.Error(err))
		return nil, fmt.Errorf("数据库错误: %w", err)
	}

	// 比较密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 生成 JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error("生成 Token 失败", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", ErrTokenGeneration, err)
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
