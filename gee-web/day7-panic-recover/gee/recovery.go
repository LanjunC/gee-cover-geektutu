package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])	//唯一疑惑的一点是为什么堆栈信息能追溯到main.go

	var str strings.Builder
	str.WriteString(message + "\nTraceBack:")
	for _,pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s\n\n", trace(fmt.Sprintf("%s",err)))
				c.Fail(http.StatusInternalServerError,"Internal Server Error")
			}
		}()
		c.Next()
	}
}
