# gin-screw
Secondary packaging based on gin framework

Make gin easier to use

Do you need to customize the validator template?

Do you need a faster way to register routes?

Do you need to use JWT?

Do you need to use zapLog?

I have all these

Make development faster!!!

### Install
1、go mod init

2、your project main.go import "github.com/gin-screw/gin-screw/ginx" and "github.com/gin-gonic/gin"

3、go mod tidy

### Use
```golang
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/miajio/gin-screw/ginx"
)

type testRouter struct{}

func (t *testRouter) Execute(c *gin.Engine) {
	c.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
}

var TestRouter ginx.Router = (*testRouter)(nil)

func main() {
	ginx.Init(gin.New())
	ginx.AddRouters(
		TestRouter,
	)
	ginx.Execute()
	ginx.Engine().Run(":8088")
}
```