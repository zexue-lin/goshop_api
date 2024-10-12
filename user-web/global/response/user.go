package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02 15:04:05"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	Nickname string `json:"name"`
	//Birthday string `json:"birthday"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}

/*
这里对时间处理的格式
主要是yapi接口文档里面 返回的birthday

为每一个特殊的返回自定义一个struct
*/
