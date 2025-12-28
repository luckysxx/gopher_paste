// handler/paste.go
package handler

import (
	"errors"
	"project/common/response"
	"project/services/paste/model"
	"project/services/paste/repository"
	"project/services/paste/service"

	"github.com/gin-gonic/gin"
)

type PasteHandler struct {
	svc service.PasteService
}

func NewPasteHandler(svc service.PasteService) *PasteHandler {
	return &PasteHandler{svc: svc}
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
		response.Error(c, err)
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
		if errors.Is(err, repository.ErrPasteNotFound) {
			response.NotFound(c, "Not Found")
			return
		}
		response.Error(c, err)
		return
	}
	response.Success(c, paste)
}
