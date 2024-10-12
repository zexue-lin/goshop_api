package main

import (
	"time"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig() // 主要用它‘输出到文件配置’的功能
	cfg.OutputPaths = []string{
		"./myproject.log",
		"stderr",
		"stdout",
	}
	return cfg.Build()
}

func main() {
	//logger, _ := zap.NewProduction()
	logger, err := NewLogger()
	if err != nil {
		panic(err)
		//panic("初始化logger失败")
	}
	su := logger.Sugar()
	defer su.Sync()
	url := "https://imooc.com"
	su.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
