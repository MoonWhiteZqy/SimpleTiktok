package dao

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type LikeModel struct {
	*gorm.Model
	UserId  int64 `gorm:"column:user_id"`
	VideoId int64 `gorm:"column:video_id"`
}

type CommentModel struct {
	*gorm.Model
	UserId  int64  `gorm:"column:user_id"`
	VideoId int64  `gorm:"column:video_id"`
	Content string `gorm:"column:content"`
}

// 实现用户给视频点赞/取消点赞
//
// TODO:对MySQL操作交给MQ
func LikeAction(userId int64, videoIdStr string, like bool) error {
	var err error
	// 更新Redis内的Like
	if like {
		_, err = rdbLike.SAdd(ctx, fmt.Sprintf("%v", userId), videoIdStr).Result()
	} else {
		_, err = rdbLike.SRem(ctx, fmt.Sprintf("%v", userId), videoIdStr).Result()
	}
	if err != nil {
		return err
	}

	// 确保被点赞的视频存在
	var video FileModel
	DB.Where("id = ?", videoIdStr).Find(&video)
	if video.Model == nil {
		return errors.New("video not exist")
	}

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		return err
	}

	// 更新数据库Like
	likeModel := LikeModel{UserId: userId, VideoId: videoId}
	return DB.FirstOrCreate(&likeModel).Error
}

// 获取用户点赞过的所有视频
func LikeList(userId int64, userHost string) ([]Video, error) {
	// 从Redis读取点赞过的视频
	videoIdStrs, err := rdbLike.SMembers(ctx, fmt.Sprintf("%v", userId)).Result()
	if err != nil {
		return make([]Video, 0), err
	}

	// 填充所有视频数据
	var res []Video
	for _, videoIdStr := range videoIdStrs {
		videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
		if err != nil {
			continue
		}
		// 从数据库获取视频数据
		video, err := getVideoById(videoId, userId, userHost)
		if err != nil {
			break
		}
		res = append(res, video)
	}
	return res, err
}

// 用户对视频发表评论
func PublishComment(userId int64, videoId int64, commentText string) (Comment, error) {
	// 把评论写入数据库
	comment := CommentModel{
		UserId:  userId,
		VideoId: videoId,
		Content: commentText,
	}
	err := DB.Create(&comment).Error
	if err != nil {
		return Comment{}, err
	}

	// 获取发布评论者的信息
	user, err := getUserById(userId, userId)
	if err != nil {
		return Comment{}, errors.New("err when getting user")
	}

	// 获取评论的日期
	month, day := comment.CreatedAt.Month(), comment.CreatedAt.Day()
	return Comment{
		CommentId:  int64(comment.ID),
		VideoId:    videoId,
		User:       user,
		Content:    comment.Content,
		CreateDate: fmt.Sprintf("%02d-%02d", month, day),
	}, nil
}

// 用户删除发表的评论
func DeleteComment(userId int64, commentId int64) error {
	return DB.Model(CommentModel{}).Where("user_id = ?", userId).Delete(&CommentModel{}, commentId).Error
}

// 获取视频的所有评论,按发布时间倒序
func GetCommentDescByPublishTime(videoId int64, jwtUserId int64) ([]Comment, error) {
	var commentLogs []CommentModel
	err := DB.Where("video_id = ?", videoId).Order("created_at desc").Find(&commentLogs).Error
	if err != nil {
		return make([]Comment, 0), err
	}

	var res []Comment
	for _, commentLog := range commentLogs {
		// 获取评论发布者的信息
		user, err := getUserById(commentLog.UserId, jwtUserId)
		if err != nil {
			continue
		}

		month, day := commentLog.CreatedAt.Month(), commentLog.CreatedAt.Day()
		res = append(res, Comment{
			CommentId:  int64(commentLog.ID),
			VideoId:    videoId,
			User:       user,
			Content:    commentLog.Content,
			CreateDate: fmt.Sprintf("%02d-%02d", month, day),
		})
	}
	return res, nil
}

// 获取所有Like数据,交给Redis
func getLikeLogs() []LikeModel {
	var likes []LikeModel
	DB.Find(&likes)
	return likes
}

// 获取所有 关注 关系,交给Redis
func getFollowLogs() []FollowModel {
	var follows []FollowModel
	DB.Find(&follows)
	return follows
}
