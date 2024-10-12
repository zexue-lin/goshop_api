package main

import (
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction() // 生成坏境的用法
	//logger, _ := zap.NewDevelopment() // 开发坏境
	defer logger.Sync() // 刷新缓存
	url := "https://imooc.com"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second, // 如果获取URL失败，程序会等待1秒再进行3此尝试
	) // 类似于Map，比较方便
	sugar.Infof("Failed to fetch URL: %s", url) // 类似与Printf

	// 直接使用logger
	// logger, _ := zap.NewProduction()
	// //logger, _ := zap.NewDevelopment()
	// defer logger.Sync()
	// url := "https://imooc.com"
	// logger.Info("failed to fetch URL",
	// 	zap.String("url", url),
	// 	zap.Int("nums", 3),
	// )

}

/*
Zap提供了两种类型的日志记录器—Sugared Logger和Logger。

这里还可以直接使用logger
它甚至比SugaredLogger更快，因为它诗经指明了类型，不会启用go语言的反射机制，启动了话相对会慢一些

绝大部分情况下使用 sugared logger 性能已经很好了，比其他日志库还要高性能
*/
