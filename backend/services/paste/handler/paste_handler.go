// handler/paste.go
package handler

import (
	"errors"
	"project/common/dberr"
	"project/common/errs"
	"project/common/response"
	"project/services/paste/model"
	"project/services/paste/service"

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

	paste, err := h.svc.Create(c.Request.Context(), &req)
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

// @Summary      获取代码片段
// @Tags         pastes
// @Produce      json
// @Param        id   path      string  true  "短链接ID"
// @Success      200  {object}  model.PasteResponse
// @Router       /pastes/{id} [get]
func (h *PasteHandler) Get(c *gin.Context) {
	shortLink := c.Param("id")

	paste, err := h.svc.GetByShortLink(c.Request.Context(), shortLink)
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
