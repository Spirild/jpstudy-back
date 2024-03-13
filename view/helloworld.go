package view

// 后续，这里会出现大量的view层函数，在httpserver中被引入并外显
// 如何对这些大量的view函数进行拆分，不是现在该考虑的事情了

import (
	"github.com/gin-gonic/gin"
)

func HelloWorld(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
