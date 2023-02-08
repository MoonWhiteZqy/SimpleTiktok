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
	DB.AutoMigrate(FollowModel{})

	// 点赞记录
	rdbLike = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   1,
	})
	// 被关注者的关注者记录
	rdbMasterFollowerDB = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   2,
	})
	// 用户关注的所有人记录
	rdbFollowerMasterDB = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   3,
	})

	// 读取MySQL数据到Redis
	rdbInit()
}
