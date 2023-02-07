package dao

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdbLike *redis.Client
var rdbMasterFollowerDB *redis.Client // key为Master
var rdbFollowerMasterDB *redis.Client // key为Follower
