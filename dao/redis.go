package dao

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdbLike *redis.Client
