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

#### ginx

```golang
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/miajio/gin-screw/pkg/ginx"
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

#### log

```golang
package main

import "github.com/miajio/gin-screw/pkg/log"

func main() {
	lo := map[string]log.Level{
		"debug.log": log.DebugLevel,
		"info.log":  log.InfoLevel,
		"error.log": log.ErrorLevel,
	}
	log.Init("./log", 256, 10, 7, false, lo)
	log.GetLogger().Info("hello")
}
```

#### jwt
```golang
package main

import (
	"fmt"
	"time"

	"github.com/miajio/gin-screw/pkg/jwt"
)

func main() {
	var params = map[string]string{}
	params["account"] = "miajio"
	params["userName"] = "admin"
	params["unitId"] = "1"
	val, err := jwt.EncryptionToken(params, "test", time.Hour*5)
	if err != nil {
		fmt.Printf("encryption token error: %v", err)
		return
	}
	fmt.Println(val)

	v2, err := jwt.DecryptionToken(val, "test")
	if err != nil {
		fmt.Printf("eecryption token error: %v", err)
		return
	}
	fmt.Println(v2)
}
```