package ginx

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// ginx
type ginx struct {
	engine  *gin.Engine // gin engine
	routers []Router    // routers
}

var (
	this *ginx
	mu   sync.Mutex
)

// Init singleton init ginx
func Init(engine *gin.Engine) *ginx {
	if this == nil {
		mu.Lock()
		defer mu.Unlock()
		if this == nil {
			this = &ginx{
				engine:  engine,
				routers: make([]Router, 0),
			}
		}
	}
	return this
}

// check ginx init
func check() {
	if this == nil {
		panic("ginx not init, please call Init()")
	}
	if this.engine == nil {
		panic("gin engine is nil; please call Init()")
	}
}

// AddRouter add router slice
func AddRouters(routers ...Router) {
	check()
	this.routers = append(this.routers, routers...)
}

// Execute execute router
func Execute() {
	check()
	for _, router := range this.routers {
		router.Execute(this.engine)
	}
}

// Engine get gin engine
func Engine() *gin.Engine {
	check()
	return this.engine
}
