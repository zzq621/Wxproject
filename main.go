package main

import (
	"WxProject/config"
	"WxProject/routers"
	"github.com/gin-gonic/gin"
)

func loadGin() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	// 使用中间件
	//r.Use(middleware.LoggerToFile())
	// 注册路由
	routers.LoadRouters(r)
	//dao.LoadDatabase()
	return r
}

func main() {
	r := loadGin()
	// logo
	//xstring.GenLogoAscii("GPT-BOT", "green")
	// 启动gin
	_ = r.Run(":" + config.GetSystemConf().Port)
	// 关闭后置操作
	//dao.CloseDb()
}
