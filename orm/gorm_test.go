package orm

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func initGorm() (cd *sql.DB, v *gorm.DB) {
	dsn := "root:adnu211nd1@tcp(192.168.56.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
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
