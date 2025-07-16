package v1

import (
	"github.com/gin-gonic/gin"
)

type AuthHandlers interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetMe(c *gin.Context)
}
