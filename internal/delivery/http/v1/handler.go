package v1

import (
	"github.com/begenov/courses-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    service.Service
	responseCh chan []byte
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service:    *service,
		responseCh: make(chan []byte),
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		go h.consumeResponseMessages()
		h.initCoursesRoutes(v1)

	}
}
