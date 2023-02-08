package dao

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdbLike *redis.Client
var rdbMasterFollowerDB *redis.Client // key为Master
var rdbFollowerMasterDB *redis.Client // key为Follower

// 根据MySQL初始化Redis数据
func rdbInit() {
	rdbLikeInit()
	rdbFollowInit()
}

// 初始化点赞的Redis表
func rdbLikeInit() {
	mySQLLikes := getLikeLogs()
	for _, mySQLLike := range mySQLLikes {
		rdbLike.SAdd(ctx, i64ToStr(mySQLLike.UserId), i64ToStr(mySQLLike.VideoId))
	}
}

// 初始化关注列表
func rdbFollowInit() {
	mySQLFollows := getFollowLogs()
	for _, mySQLFollow := range mySQLFollows {
		rdbFollowerMasterDB.SAdd(ctx, i64ToStr(mySQLFollow.FollowerId), i64ToStr(mySQLFollow.MasterId))
		rdbMasterFollowerDB.SAdd(ctx, i64ToStr(mySQLFollow.MasterId), i64ToStr(mySQLFollow.FollowerId))
	}
}
