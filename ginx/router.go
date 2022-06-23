package ginx

import "github.com/gin-gonic/gin"

// Router ginx
type Router interface {
	Execute(c *gin.Engine) // execute router
}
