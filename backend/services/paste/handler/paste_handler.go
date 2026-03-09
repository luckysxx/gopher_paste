// handler/paste.go
package handler

import (
	"errors"
	"project/common/dberr"
	"project/common/errs"
	"project/common/response"
	"project/services/paste/model"
	"project/services/paste/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PasteHandler struct {
	svc    service.PasteService
	logger *zap.Logger
}

func NewPasteHandler(svc service.PasteService, logger *zap.Logger) *PasteHandler {
	return &PasteHandler{svc: svc, logger: logger}
}

// @Summary      创建代码片段
// @Tags         pastes
// @Accept       json
// @Produce      json
// @Param        request body model.CreatePasteRequest true "请求参数"
// @Success      200  {object}  model.PasteResponse
// @Router       /pastes [post]
func (h *PasteHandler) Create(c *gin.Context) {
	var req model.CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	paste, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("创建 paste 失败", zap.Error(err))

		// 短链接生成失败需要特殊消息
		if errors.Is(err, service.ErrShortLinkGeneration) {
			response.Error(c, errs.New(errs.ServerErr, "系统繁忙，请稍后再试", err))
			return
		}

		// 其他错误统一转换
		response.Error(c, errs.ConvertToCustomError(err))
		return
	}
	response.Success(c, paste)
}

// @Summary      获取我的代码片段列表
// @Tags         pastes
// @Produce      json
// @Success      200  {array}  model.PasteResponse
// @Router       /me/pastes [get]
func (h *PasteHandler) ListMine(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	list, err := h.svc.ListMine(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("查询用户 paste 列表失败", zap.Error(err))
		response.Error(c, errs.ConvertToCustomError(err))
		return
	}

	response.Success(c, list)
}

// @Summary      获取代码片段
// @Tags         pastes
// @Produce      json
// @Param        id   path      string  true  "片段ID"
// @Success      200  {object}  model.PasteResponse
// @Router       /pastes/{id} [get]
func (h *PasteHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "无效的片段ID")
		return
	}

	paste, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("查询 paste 失败", zap.Error(err))

		// NotFound 错误需要特殊的 HTTP 状态码
		if dberr.IsNotFoundError(err) {
			response.NotFound(c, "Paste 不存在")
			return
		}

		// 其他错误统一转换
		response.Error(c, errs.ConvertToCustomError(err))
		return
	}
	response.Success(c, paste)
}

// @Summary      更新代码片段
// @Tags         pastes
// @Accept       json
// @Produce      json
// @Param        id path string true "片段ID"
// @Param        request body model.UpdatePasteRequest true "请求参数"
// @Success      200  {object}  model.PasteResponse
// @Router       /pastes/{id} [put]
func (h *PasteHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "无效的片段ID")
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}

	var req model.UpdatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	paste, err := h.svc.Update(c.Request.Context(), userID, id, &req)
	if err != nil {
		h.logger.Error("更新 paste 失败", zap.Error(err))

		if errors.Is(err, service.ErrForbidden) {
			response.Forbidden(c, "无权限操作该代码片段")
			return
		}

		response.Error(c, errs.ConvertToCustomError(err))
		return
	}

	response.Success(c, paste)
}

func getUserID(c *gin.Context) (int64, bool) {
	v, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	userID, ok := v.(int64)
	if !ok {
		return 0, false
	}

	return userID, true
}
