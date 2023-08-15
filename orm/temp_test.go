package orm

import (
	"testing"
	"time"
)

type RealPeopleCert struct {
	Id            int32      `gorm:"column:id" json:"id"`
	Uid           int64      `gorm:"column:uid" json:"uid"`
	SamePerc      int32      `gorm:"column:same_perc" json:"same_perc"`
	CertImg       string     `gorm:"column:cert_img" json:"cert_img"`
	IpLoc         string     `gorm:"column:ip_loc" json:"ip_loc"`
	IsFirstUpload bool       `gorm:"column:is_first_upload" json:"is_first_upload"`
	AuditUid      int32      `gorm:"column:audit_uid" json:"audit_uid"`
	AuditStatus   int32      `gorm:"column:audit_status" json:"audit_status"`
	UploadAt      time.Time  `gorm:"column:upload_at" json:"upload_at"`
	AuditAt       *time.Time `gorm:"column:audit_at" json:"audit_at"`
}

type RealPeopleCertT struct {
	RealPeopleCert
	Nick      string `gorm:"column:nick" json:"nick"`
	Sex       int32  `gorm:"column:sex" json:"sex"`
	Avatar    string `gorm:"column:avatar" json:"avatar"`
	AuditUser string `gorm:"column:audit_user" json:"audit_user"`
}

func (RealPeopleCert) TableName() string {
	return "real_people_cert"
}

func TestX(t *testing.T) {
	cc := mustLoadConf()
	db, v := initGorm(cc.GormTestdb.Dsn)
	defer db.Close()

	sql := `select r.*, u.nickname nick, u.sex, u.avatar
			from melon_user.real_people_cert r
					 left join melon_user.user_info u on r.uid = u.userid
			where 1 = 1
			  and audit_status = 2
			  and u.sex = 1
			  and r.is_first_upload = 1
			order by r.upload_at
			limit 1;`
	var list []RealPeopleCertT
	err := v.Raw(sql).Scan(&list).Error
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	t.Logf("LIST:%+v", list[0])
}
