package dao

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type FollowModel struct {
	*gorm.Model
	MasterId   int64 `gorm:"column:master_id"`   // 被关注者的id
	FollowerId int64 `gorm:"column:follower_id"` // 点击关注的用户的id
}

func i64ToStr(id int64) string {
	return fmt.Sprintf("%v", id)
}

// 关注用户
func FollowMaster(userId, masterId int64) error {
	// 添加到redis表中
	n1, err := rdbMasterFollowerDB.SAdd(ctx, i64ToStr(masterId), i64ToStr(userId)).Result()
	if err != nil {
		return fmt.Errorf("err in redis master table: %v", err)
	}
	n2, err := rdbFollowerMasterDB.SAdd(ctx, i64ToStr(userId), i64ToStr(masterId)).Result()
	if err != nil {
		return fmt.Errorf("err in redis follower table: %v", err)
	}
	fmt.Println("affected:", n1, n2)

	// 添加到MySQL中
	err = DB.FirstOrCreate(&FollowModel{MasterId: masterId, FollowerId: userId}).Error
	return err
}

// 取关用户
func UnFollowMaster(userId, masterId int64) error {
	// 从Redis表中删除记录
	n1, err := rdbMasterFollowerDB.SRem(ctx, i64ToStr(masterId), i64ToStr(userId)).Result()
	if err != nil {
		return fmt.Errorf("err in redis master table: %v", err)
	}
	n2, err := rdbFollowerMasterDB.SRem(ctx, i64ToStr(userId), i64ToStr(masterId)).Result()
	if err != nil {
		return fmt.Errorf("err in redis follower table: %v", err)
	}
	fmt.Println("affected:", n1, n2)

	// 从MySQL中删除
	err = DB.Where("master_id = ?", masterId).Where("follower_id = ?", userId).Delete(&FollowModel{}).Error
	return err
}

// 获取 用户 在对应表中的所有用户信息
func getUserCorrespondList(userId string, rdb *redis.Client) ([]User, error) {
	// 获取表中存储的 所有用户id
	targetIdStrs, err := rdb.SMembers(ctx, userId).Result()
	if err != nil {
		return make([]User, 0), err
	}

	var userList []User
	// 通过id获取用户信息
	for _, targetIdStr := range targetIdStrs {
		user, err := getUserByIdStr(targetIdStr)
		if err != nil {
			continue
		}
		userList = append(userList, user)
	}
	return userList, err
}

// 获取 当前用户 关注的所有用户
func GetMasterList(userId string) ([]User, error) {
	users, err := getUserCorrespondList(userId, rdbFollowerMasterDB)
	// 由于获取的所有用户都是关注的用户,因此把IsFollow都置为true
	for i := range users {
		users[i].IsFollow = true
	}
	return users, err
}

// 获取 当前用户 关注的所有用户
func GetFollowerList(masterId string) ([]User, error) {
	users, err := getUserCorrespondList(masterId, rdbMasterFollowerDB)
	return users, err
}
