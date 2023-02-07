package dao

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"
)

type FileModel struct {
	*gorm.Model
	AuthorId int64  `gorm:"column:author_id"`
	Title    string `gorm:"column:title"`
	FileName string `gorm:"column:title"`
}

// 将上传的文件保存到文件夹,并添加到数据库
//
// TODO:生成封面、上传记录唯一性
func SavePOSTFile(content []byte, path, fileName, title string, userId int64) error {
	var err error
	// 创建用户专属的文件夹
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	// 把上传的文件写到磁盘
	err = os.WriteFile(path+fileName, content, 0777)
	if err != nil {
		return err
	}

	// 在数据库添加文件记录
	DB.Create(&FileModel{AuthorId: userId, Title: title, FileName: fileName})
	return DB.Error
}

// 获取对应用户id上传的视频列表
//
// TODO: 完善Video和Author信息
func GetVideoOfUser(userId int64) ([]Video, error) {
	var videos []FileModel
	var res []Video
	DB.Where("author_id = ?", userId).Find(&videos)
	for _, v := range videos {
		author, err := getUserById(userId)
		if err != nil {
			fmt.Println("err when read author", err)
			continue
		}
		res = append(res, Video{
			VideoId:       int64(v.ID),
			Author:        author,
			PlayUrl:       "",
			CoverUrl:      "",
			FavoriteCount: 0,
			CommentCount:  0,
			Title:         v.Title,
		})
	}
	return res, DB.Error
}

// 获取Feed
//
// TODO:完善视频内容Url
func GetFeed(latestTime int64) ([]Video, int64, error) {
	var videos []FileModel
	var res []Video
	nextTime := int64(-1)
	// 获取视频文件
	DB.Where("created_at < ?", time.Unix(latestTime, 0)).Order("created_at desc").Limit(30).Find(&videos)
	for _, v := range videos {
		author, err := getUserById(v.AuthorId)
		if err != nil {
			fmt.Println("err when read author", err)
			continue
		}
		res = append(res, Video{
			VideoId:       int64(v.ID),
			Author:        author,
			PlayUrl:       "",
			CoverUrl:      "",
			FavoriteCount: 0,
			CommentCount:  0,
			Title:         v.Title,
		})
		nextTime = v.CreatedAt.Unix()
	}
	return res, nextTime, DB.Error
}

// 在数据库中查找对应Id的视频
func getVideoById(videoId int64) (Video, error) {
	// 根据id获取视频
	var videoFIle FileModel
	err := DB.Where("id = ?", videoId).Find(&videoFIle).Error
	if err != nil {
		return Video{}, errors.New("err when getting video")
	}

	// 获取视频上传者信息
	author, err := getUserById(videoFIle.AuthorId)
	if err != nil {
		return Video{}, errors.New("err when getting author")
	}

	// 填充返回内容
	video := Video{
		VideoId:       videoId,
		Author:        author,
		PlayUrl:       "",
		CoverUrl:      "",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         videoFIle.Title,
	}
	return video, nil
}
