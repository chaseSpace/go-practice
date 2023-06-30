package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go_accost/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var M *gorm.DB
var R *redis.Client

func Init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/x_accost?charset=utf8mb4&parseTime=True&loc=Local",
		config.V.Mysql.Username, config.V.Mysql.Password,
		config.V.Mysql.Host, config.V.Mysql.Port)
	M, err = gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// REDIS
	R = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%d`, config.V.Redis.Host, config.V.Redis.Port),
		Password: config.V.Redis.Password,
		DB:       config.V.Redis.DB,
	})
	err = R.Ping(context.TODO()).Err()
	if err != nil {
		panic(err)
	}
}
