/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:07
 * Filename      : user.go
 * Description : 用户基础存储数据
 * *******************************************************/
package data

import (
	"basic/ssdb/gossdb"
	"basic/utils"
	"sync"

	"github.com/golang/glog"
)

var tex sync.Mutex

// 生成玩家唯一id
func GenerateUserid() (string, error) {
	tex.Lock()
	value, err := gossdb.C().Get(KEY_LAST_USER_ID)
	userid := string(value)
	if userid == "" {
		userid = "60001"
	}

	err = gossdb.C().Set(KEY_LAST_USER_ID, utils.StringAdd(userid))

	tex.Unlock()
	return userid, err
}
func GetMultiUser(userids []string) []*User {
	users := make([]*User, 0, len(userids))
	for _, v := range userids {
		user := &User{Userid: v}
		if err := user.Get(); err != nil {
			glog.Errorln(err)
		}
		users = append(users, user)
	}
	return users
}
func (this *User) Get() error {
	return gossdb.C().GetObject(KEY_USER+this.Userid, this)
}
func (this *User) Save() error {
	return gossdb.C().PutObject(KEY_USER+this.Userid, this)
}

func (this *User) ExistPhone(phone string) bool {
	value, err := gossdb.C().Hget(KEY_USER+this.Userid, "Phone")
	if err != nil {
		return false
	}
	return string(value) == phone
}

func (this *User) ExistNickname(nickname string) bool {
	value, err := gossdb.C().Hget(KEY_USER+this.Userid, "Nickname")
	if err != nil {
		return false
	}
	return string(value) == nickname
}

func (this *User) UpdateSex() error {
	return gossdb.C().Hset(KEY_USER+this.Userid, "Sex", this.Sex)
}

func (this *User) UpdateNickname() error {
	return gossdb.C().Hset(KEY_USER+this.Userid, "Nickname", this.Nickname)
}
func (this *User) GetByPhone() string {
	value, err := gossdb.C().Hget(KEY_USER+this.Userid, "Phone")
	if err != nil {
		return ""
	}
	return string(value)
}
func (this *User) UpdatePWD(pwd string) error {
	auth, err := gossdb.C().Hget(KEY_USER+this.Userid, "Auth")
	passwd := utils.Md5(pwd + string(auth))
	err = gossdb.C().Hset(KEY_USER+this.Userid, "Passwd", passwd)
	return err
}

//  用户登陆密码验证
func (this *User) PWDIsOK(pwd string) bool {
	value, err := gossdb.C().MultiHget(KEY_USER+this.Userid, "Pwd", "Auth")
	if err != nil {
		glog.Infoln(err, string(value["Pwd"]))
		return false
	}
	//  密码正确
	if utils.Md5(pwd+string(value["Auth"])) == string(value["Pwd"]) {
		return true
	}
	glog.Infoln(err, string(value["Pwd"]), string(value["Auth"]))
	return false
}

type User struct {
	Userid        string // 用户id
	Nickname      string // 用户昵称
	Sex           uint32 // 用户性别,男1 女2 非男非女3
	Sign          string // 用户签名
	Email         string // 绑定的邮箱地址
	Phone         string // 绑定的手机号码
	Auth          string // 密码验证码
	Pwd           string // MD5密码
	Birth         uint32 // 用户生日日期
	Create_ip     uint32 // 注册账户时的IP地址
	Create_time   uint32 // 注册时间
	Coin          uint32 // 金币
	Exp           uint32 // 经验
	Diamond       uint32 // 钻石
	Ticket        uint32 //入场券
	Exchange      uint32 //兑换券
	Terminal      string // 终端类型名字
	Status        uint32 // 正常1  锁定2  黑名单3
	Address       string //物理地址
	Photo         string //头像
	Qq_uid        string //
	Wechat_uid    string
	Microblog_uid string
	Vip           uint32
	Win           uint32
	Lost          uint32
	Ping          uint32
	Platform      uint32
	VipExpire     uint32
	ChenmiTime    uint32 // 防沉迷限制
	Chenmi        int32  // 防沉迷限制
}