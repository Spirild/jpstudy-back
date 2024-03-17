package main

// 该部分引入具体类以使相关的init函数被调用, 注册

import (
	_ "translasan-lite/db"
	_ "translasan-lite/view"
)
