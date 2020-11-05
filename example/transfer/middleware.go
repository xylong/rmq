package transfer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckError 错误检查
func CheckError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s%s", msg, err.Error()))
	}
}

// ErrorMiddleware 错误处理
func ErrorMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
			}
		}()
		context.Next()
	}
}
