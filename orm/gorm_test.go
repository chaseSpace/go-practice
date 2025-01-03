package orm

import (
	"database/sql"
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"testing"
)

type Conf struct {
	GormTestdb struct {
		Dsn string `json:"dsn"`
	} `json:"gorm_testdb"`
}

func mustLoadConf() *Conf {
	b, _ := os.ReadFile("../_ignore/db.json")
	if b == nil {
		panic("no db.json found")
	}
	cc := new(Conf)
	_ = json.Unmarshal(b, cc)
	return cc
}

func initGorm(__dsn ...string) (cd *sql.DB, v *gorm.DB) {
	dsn := "root:adnu211nd1@tcp(192.168.56.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	if len(__dsn) > 0 {
		dsn = __dsn[0]
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}
	d, _ := db.DB()
	return d, db
}

// 支持指针作为参数
func TestGormSupportPtrArgs(t *testing.T) {
	db, v := initGorm()
	defer db.Close()

	x := 1
	str := "x"
	err := v.Exec(`insert into t1 (id, tt, str) value (1, ?, ?)`, &x, &str).Error
	t.Error(err)
}
