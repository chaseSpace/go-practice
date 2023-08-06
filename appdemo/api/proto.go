package api

import (
	"fmt"
	"go_accost/model"
	"go_accost/util"
	"time"
)

type BaseRsp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// -------------------------------------------------

// GetUserInfoReq 获取单个用户信息：在线状态
type GetUserInfoReq struct {
	Uid int32 `json:"uid"`
}
type GetUserInfoRes struct {
	BaseRsp
	IsOnline bool `json:"is_online"`
}

// GetTwoUserStatReq 获取两个用户之间的拉黑情况、亲密度
type GetTwoUserStatReq struct {
	Uid1 int32 `json:"uid1"`
	Uid2 int32 `json:"uid2"`
}
type GetTwoUserStatRes struct {
	BaseRsp
	Block    int8  `json:"block"`     // 0-互未拉黑，1-被一方拉黑
	CloseVal int32 `json:"close_val"` // 亲密度
}

// GetRegisterUsersReq 获取注册用户数据（筛选注册时间>=${RegStartTime} && 按注册时间升序 && 取 ${Ps} 条数）
type GetRegisterUsersReq struct {
	RegStartTime string `json:"reg_start_time"` // yyyy-mm-dd hh:MM:ss
	Ps           int    `json:"ps"`             // size
}
type GetRegisterUsersRes struct {
	BaseRsp
	List []*RegisterUser `json:"list"` // 按 RegTime 升序
}

type RegisterUser struct {
	Uid            int32        `json:"uid"`
	Name           string       `json:"name"` // 昵称
	Gender         model.Gender `json:"gender"`
	Birthday       string       `json:"birthday"` // yyyy-mm-dd hh:MM:ss
	_Birthday      time.Time
	RegTime        string `json:"reg_time"` // yyyy-mm-dd hh:MM:ss
	_RegTime       time.Time
	CertifiedTime  string `json:"certified_time"` // 真人认证时间 yyyy-mm-dd hh:MM:ss, 无则空
	_CertifiedTime *time.Time

	Income int32 `json:"income"`
}

func (r *RegisterUser) Check() (err error) {
	r._Birthday, err = time.Parse(util.Datetime, r.Birthday)
	if err != nil {
		return fmt.Errorf("birthday err: %v", err)
	}
	r._RegTime, err = time.Parse(util.Datetime, r.RegTime)
	if err != nil {
		return fmt.Errorf("RegTime err: %v", err)
	}
	if r.CertifiedTime != "" {
		ct, err := time.Parse(util.Datetime, r.CertifiedTime)
		if err != nil {
			return fmt.Errorf("CertifiedTime err: %v", err)
		}
		r._CertifiedTime = &ct
	}
	return nil
}

func (r *RegisterUser) ToUser() *model.User {
	return &model.User{
		Id:                         0,
		Uid:                        r.Uid,
		Name:                       r.Name,
		Gender:                     r.Gender,
		Birthday:                   r._Birthday,
		RegisTime:                  r._RegTime,
		CertifiedTime:              r._CertifiedTime,
		Income:                     r.Income,
		BeAccostFlow:               0,
		AccostFlow:                 0,
		BeAccostFlowUpperLimit:     0,
		AccostFlowUpperLimit:       0,
		BeAccostFlowStarUpperLimit: 0,
		AccostFlowStarUpperLimit:   0,
		CateId:                     0,
		Quality:                    "",
	}
}

// GetUserDayIncomeReq 指定日期范围内的全部用户的收益按天汇总数据
type GetUserDayIncomeReq struct {
	DayStart string `json:"day_start"` // yyyy-mm-dd
	DayEnd   string `json:"day_end"`   // yyyy-mm-dd
}
type GetUserDayIncomeRes struct {
	BaseRsp
	List []*DayIncome `json:"list"`
}

type DayIncome struct {
	Day    string `json:"day"` // yyyy-mm-dd
	Uid    int32  `json:"uid"`
	Income int32  `json:"income"` // 收益：元
}
