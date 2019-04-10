package controller

import (
	"wScheduler/common"
	"wScheduler/entity"

	"github.com/go-redis/redis"
)

type MailboxController struct {
	BaseController
}

func (this *MailboxController) Active() {
	res := new(common.Result)
	defer func() {
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}()

	mailbox := this.GetString("Mailbox")
	key := this.GetString("Key")

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		res.Message = err.Error()
		return
	}

	localKey, err := client.HGet("mail:"+mailbox, "key").Result()
	if err != nil {
		res.Message = err.Error()
		return
	}

	if localKey != key {
		res.Message += " 未通过验证"
		return
	} else {
		res.Message += " 通过验证"
	}

	var m []interface{}
	m, err = client.HMGet("mail:"+mailbox, "username", "password").Result()
	if err != nil {
		res.Message += " 注册失败"
		return
	}
	userInfo := new(entity.UserInfo)
	userInfo.Username = m[0].(string)
	userInfo.Password = m[1].(string)
	userInfo.Mailbox = mailbox
	err = userInfo.SaveUserInfo()
	if err != nil {
		res.Message += "注册失败"
		return
	}
	res.Message += " 注册成功"
	res.Success = true
}
