// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. Created at 2024-05-09 18:37:52
// =================================================================================

package entity

import (
	"time"
)

// Users is the golang structure for table users.
type Users struct {
	UserId           int       `json:"userId"           orm:"user_id"           ` // 用户ID，自增主键
	Username         string    `json:"username"         orm:"username"          ` // 用户名，最大长度为50个字符，不为空
	Email            string    `json:"email"            orm:"email"             ` // 电子邮件地址，最大长度为100个字符，不为空，唯一索引
	Password         string    `json:"password"         orm:"password"          ` // 密码，最大长度为255个字符，不为空
	Gender           bool      `json:"gender"           orm:"gender"            ` // 性别，使用 BIT(1) 类型表示，0 表示女性，1 表示男性，不为空
	Money            float64   `json:"money"            orm:"money"             ` // money，浮点数
	RegistrationDate time.Time `json:"registrationDate" orm:"registration_date" ` // 注册日期，默认为当前时间戳
}
