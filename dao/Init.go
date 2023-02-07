package dao

import (
	"fmt"
	"simpleTiktok/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.DSN))
	if err != nil {
		fmt.Println(err)
	}
	DB.AutoMigrate(UserModel{})
	DB.AutoMigrate(FileModel{})
	DB.AutoMigrate(LikeModel{})
	DB.AutoMigrate(CommentModel{})

	// 点赞记录
	rdbLike = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   1,
	})
}
