package initialize

import (
	"fmt"
	"goshop_api/user-web/global"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("GOSHOP_DEBUG")

	configFilePrefix := "config" // 自定义文件前缀
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()

	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 这个文件如何在其他文件中使用 - 要把它配置成全局变量，那么全局变量又要放在global文件夹下
	// serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置信息:%v", global.ServerConfig)

	// viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件产生变化:%s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)

		zap.S().Infof("配置信息:%v", global.ServerConfig)
	})

	time.Sleep(time.Second * 300)
}
