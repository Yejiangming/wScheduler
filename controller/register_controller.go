package controller

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
	"wScheduler/common"
	"wScheduler/entity"

	"github.com/astaxie/beego/config"
	"github.com/go-redis/redis"
)

type RegisterController struct {
	BaseController
}

func (this *RegisterController) ToRegister() {
	this.TplName = "register/register.html"
}

func (this *RegisterController) Register() {
	res := new(common.Result)

	defer func() {
		this.Data["res"] = res
		this.TplName = "jobinfo/res.html"
	}()

	username := this.GetString("Username")
	password := this.GetString("Password")
	mailbox := this.GetString("Mailbox")

	if username == "" || password == "" || mailbox == "" {
		res.Message = " 有未填字段"
		return
	}

	userInfo := new(entity.UserInfo)
	userInfo.Username = username
	userInfo.GetUserInfo()
	if userInfo.Id != 0 {
		res.Message = " 用户名已经被注册"
		return
	}

	userInfo.Id = 0
	userInfo.Username = ""
	userInfo.Password = ""
	userInfo.Mailbox = mailbox
	userInfo.GetUserInfo()
	if userInfo.Id != 0 {
		res.Message = " 邮箱已被注册"
		return
	}

	rand.Seed(time.Now().Unix())
	h := md5.New()
	h.Write([]byte(mailbox + password))
	key := hex.EncodeToString(h.Sum(nil))

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
	m := make(map[string]interface{})
	m["username"] = username
	m["password"] = password
	m["mailbox"] = mailbox
	m["key"] = key
	err = client.HMSet("mail:"+mailbox, m).Err()
	if err != nil {
		res.Message = err.Error()
		return
	}
	err = client.Expire("mail:"+mailbox, 2*time.Minute).Err()
	if err != nil {
		res.Message = err.Error()
		return
	}

	conf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		res.Message = err.Error()
		return
	}
	servername := conf.String(conf.String("runmode") + "::servername")
	msg := "<html>\r\n" +
		"<a href=" + servername + ":8888/mailbox/active?Mailbox=" + mailbox + "&Key=" + key + ">key</a>\r\n" +
		"</html>\r\n"
	err = common.SendMail(mailbox, msg, "激活邮箱")
	if err != nil {
		res.Message += " 发送邮件失败"
		fmt.Println(err)
		return
	}

	res.Success = true
	res.Message += "发送邮件成功"
}
