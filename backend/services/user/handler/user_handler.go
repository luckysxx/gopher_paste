package handler

import (
	"project/common/response"
	"project/services/user/model"
	"project/services/user/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	Log *zap.Logger
	svc service.UserService
}

func NewUserHandler(svc service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, Log: log}
}

// @Summary      用户注册
// @Description  用户注册接口
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body model.RegisterRequest true "注册信息"
// @Success      200  {object}  model.RegisterResponse
// @Router       /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("参数绑定失败", zap.Error(err))
		response.BadRequest(c, "参数错误")
		return
	}

	user, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		h.Log.Error("用户注册失败", zap.Error(err))
		response.Error(c, err)
		return
	}
	response.Success(c, user)
}

// @Summary      用户登录
// @Description  用户登录接口
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "登录信息"
// @Success      200  {object}  model.LoginResponse
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("参数绑定失败", zap.Error(err))
		response.BadRequest(c, "参数错误")
		return
	}

	user, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		h.Log.Error("用户登录失败", zap.Error(err))
		response.Error(c, err)
		return
	}
	response.Success(c, user)
}
