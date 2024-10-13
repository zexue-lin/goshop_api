package main

import (
	"fmt"
	"goshop_api/user-web/global"
	"goshop_api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {

	// 1. 初始化logger
	initialize.InitLogger()

	// 2. 初始化配置文件
	initialize.InitConfig()

	// 3. 初始化router
	Router := initialize.Routers()

	// 4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	/*
			1. S() 可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
			2. 日志是分级别的，debug， info， warn ，error， fetal （由低到高）
		  3. S函数和L函数很有用，给我们提供了全局安全访问logger的途径
	*/
	zap.S().Debugf("启动服务器,端口:%d", global.ServerConfig.Port) // 省略代码
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
