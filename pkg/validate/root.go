package validate

import (
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	validateMap map[string]validator.Func
	mu          sync.Mutex
)

// Put put validator func
func Put(key string, value validator.Func) {
	mu.Lock()
	defer mu.Unlock()
	if validateMap == nil {
		validateMap = make(map[string]validator.Func)
	}
	validateMap[key] = value
}

// Execute execute puts validator
func Execute() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for key, val := range validateMap {
			v.RegisterValidation(key, val)
		}
	}
}
