package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// 先测试一下验证码对不对
func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%d:%s", "192.168.147.131", 6379),
	})

	value, err := rdb.Get(context.Background(), "18787636723").Result()
	if err == redis.Nil {
		fmt.Println("key 不存在")
	}
	fmt.Println(value)
	/*
		这里需要配置一下docker中 redis开机自启动，不然出错
	*/
}
