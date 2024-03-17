package view

import (
	"translasan-lite/core"
)

// module register
func init() {
	core.RegisterCompType("HttpService", (*HttpServer)(nil))
}

// type IHelloWorld interface {
// }
// 不需要被外部调用，这里就是web的入口
