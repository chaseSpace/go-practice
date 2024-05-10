// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2024-05-09 18:37:52
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta           `orm:"table:users, do:true"`
	UserId           interface{} // 用户ID，自增主键
	Username         interface{} // 用户名，最大长度为50个字符，不为空
	Email            interface{} // 电子邮件地址，最大长度为100个字符，不为空，唯一索引
	Password         interface{} // 密码，最大长度为255个字符，不为空
	Gender           interface{} // 性别，使用 BIT(1) 类型表示，0 表示女性，1 表示男性，不为空
	Money            interface{} // money，浮点数
	RegistrationDate interface{} // 注册日期，默认为当前时间戳
}
