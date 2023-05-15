package handler

import (
	"github.com/gin-gonic/gin"
)

type IHandler interface {
	Init(router *gin.RouterGroup)
}
