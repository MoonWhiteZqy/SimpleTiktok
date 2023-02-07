package service

import "simpleTiktok/dao"

type PublishService interface {
	SavePOSTFile([]byte, string, string, string, int64) error
	GetVideoOfUser(int64) ([]dao.Video, error)
	GetFeed() ([]dao.Video, error)
}
