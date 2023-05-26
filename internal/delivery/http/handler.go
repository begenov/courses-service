package http

import (
	v1 "github.com/begenov/courses-service/internal/delivery/http/v1"
	"github.com/begenov/courses-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: *service,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	h.initAPI(router)
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(&h.service)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
