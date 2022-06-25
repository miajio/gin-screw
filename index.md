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

#### validate

Built in library developed based on validator to simplify users' use of custom validators

the validate built in 
```golang
// DeSpace delete space in the val
func DeSpace(val string) string {}

// EnglishLimiter english limiter
func EnglishLimiter(fl validator.FieldLevel) bool {}

// IntegerLimiter positive integer limiter
func IntegerLimiter(fl validator.FieldLevel) bool {}

// NumberLimiter positive number limiter
func NumberLimiter(fl validator.FieldLevel) bool {}

// EqNowDayLimiter equal now day limiter
func EqNowDayLimiter(fl validator.FieldLevel) bool {}

// GtNowDayLimiter greater than now day limiter
func GtNowDayLimiter(fl validator.FieldLevel) bool {}

// LtNowDayLimiter less than now day limiter
func LtNowDayLimiter(fl validator.FieldLevel) bool {}
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

#### fileutil
Read() ([]byte, error)                          // based on reading current file data and return the file bytes

GetPath() string                                // get the file path

GetName() string                                // get the file name

GetPrefix() string                              // get the file prefix name demo: test.abc return test

GetSuffix() string                              // get the file suffix name demo: test.abc return .abc

Size() int64                                    // get the file size

IsDir() bool                                    // the file is a folder

MkdirAll(name string) (*File, error)            // based on the current folder create a new folder

Remove() error                                  // remove the current file

Rename(name string) error                       // based on the current file rename to a new file

Move(path string) error                         // based on the current file move to a new file

Paste(newpath string) error                     // based on the current file paste to a new file

copyFile(src, dest string) (w int64, err error) // private function to copy file

pathExists(path string) (bool, error)           // private function to check file path exists

GetChildren() ([]*File, error)                  // based on the current folder get all the children files

Clean()                                         // clean the file

Replace(newFile *File)                          // based on the current file replace to a new file