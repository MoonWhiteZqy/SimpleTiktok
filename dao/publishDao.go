package dao

import (
	"fmt"
	"os"
	"os/exec"
	"simpleTiktok/config"
	"time"

	"gorm.io/gorm"
)

type FileModel struct {
	*gorm.Model
	AuthorId int64  `gorm:"column:author_id"`
	Title    string `gorm:"column:title"`
	FileName string `gorm:"column:file_name"`
}

// 获取视频文件在数据库中的数据后,处理转换成[]Video
//
// userHost用来填充PlayUrl和CoverUrl
//
// TODO:完善Video点赞数量、 Video评论数量
func videoModelToVideo(models []FileModel, jwtUserId int64, userHost string) ([]Video, error) {
	res := make([]Video, 0)
	errList := make([]error, 0)
	var resErr error
	for _, model := range models {
		// 获取作者信息
		authorId := model.AuthorId
		author, err := getUserById(authorId, jwtUserId)
		if err != nil {
			errList = append(errList, err)
			continue
		}

		// 获取点赞信息
		isFavorite, _ := rdbLike.SIsMember(ctx, i64ToStr(jwtUserId), i64ToStr(int64(model.ID))).Result()

		res = append(res, Video{
			VideoId:       int64(model.ID),
			Author:        author,
			PlayUrl:       fmt.Sprintf("%v/douyin/static/%v/%v", userHost, author.UserId, model.FileName),
			CoverUrl:      fmt.Sprintf("%v/douyin/static/%v/%v", userHost, author.UserId, getCoverAddr("", model.FileName)),
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    isFavorite,
			Title:         model.Title,
		})
	}
	if len(errList) == 0 {
		resErr = nil
	} else {
		resErr = fmt.Errorf("%v", errList)
	}
	return res, resErr
}

// 生成 上传视频封面 的文件名
func getCoverAddr(path, fileName string) string {
	var idx int
	// 获取文件名后缀位置
	for i, c := range fileName {
		if c == '.' {
			idx = i
		}
	}
	if idx == 0 {
		return ""
	} else {
		return fmt.Sprintf("%v%v.jpg", path, fileName[:idx])
	}
}

// 将上传的文件保存到文件夹,并添加到数据库
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
	err = DB.Create(&FileModel{AuthorId: userId, Title: title, FileName: fileName}).Error
	if err != nil {
		return err
	}
	coverPath := getCoverAddr(path, fileName)
	fmt.Println("coverPath:", coverPath)
	if len(coverPath) > 0 {
		// var out []byte
		// Command中,空格需要用参数分割,而不能直接在string中写空格
		cmd := exec.Command("ffmpeg", "-ss", "00:00:01", "-i", fmt.Sprintf("%v%v", path, fileName), "-vframes", "1", coverPath, "-y")
		_, err = cmd.CombinedOutput()
	}
	return err
}

// 获取对应用户id上传的视频列表
func GetVideoOfUser(userId, jwtUserId int64, userHost string) ([]Video, error) {
	var videos []FileModel
	var res []Video
	var err error
	err = DB.Where("author_id = ?", userId).Find(&videos).Error
	if err != nil {
		return res, err
	}
	res, err = videoModelToVideo(videos, jwtUserId, userHost)
	return res, err
}

// 获取早于latestTime的最近数个Feed,数量由config决定
func GetFeed(latestTime, jwtUserId int64, userHost string) (res []Video, nextTime int64, err error) {
	var videos []FileModel

	// 获取视频文件信息
	err = DB.Where("created_at < ?", time.Unix(latestTime, 0)).Order("created_at desc").Limit(config.FeedLimit).Find(&videos).Error
	if err != nil {
		return
	}
	res, err = videoModelToVideo(videos, jwtUserId, userHost)

	// 更新latestTime, 获取最早时间
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].CreatedAt.Unix()
	}
	return
}

// 在数据库中查找对应Id的视频
func getVideoById(videoId, jwtUserId int64, userHost string) (Video, error) {
	// 根据id获取视频
	var videoFIle FileModel
	err := DB.Where("id = ?", videoId).Find(&videoFIle).Error
	if err != nil {
		return Video{}, fmt.Errorf("err when getting video: %v", err)
	}

	// 获取视频信息
	res, err := videoModelToVideo([]FileModel{videoFIle}, jwtUserId, userHost)
	if len(res) == 0 {
		return Video{}, fmt.Errorf("no video found")
	}
	return res[0], err
}
