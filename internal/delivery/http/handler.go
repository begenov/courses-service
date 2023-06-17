package http

import (
	"github.com/begenov/courses-service/internal/service"

	v1 "github.com/begenov/courses-service/internal/delivery/http/v1"

	_ "github.com/begenov/courses-service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
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
