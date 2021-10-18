package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"runtime/debug"
)

func ToString(i interface{}) string {
	switch t := i.(type) {
	case error:
		return t.Error()
	default:
		return t.(string)
	}
}

func PanicRecover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic: %v\n", r)
				debug.PrintStack()

				// here write the response body

				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}

func GlobalRecoverableErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
