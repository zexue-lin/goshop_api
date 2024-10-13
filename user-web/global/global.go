package global

import (
	"goshop_api/user-web/config"

	ut "github.com/go-playground/universal-translator"
)

var (
	// 定义一个全局翻译,导出去  首字母要大写
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
