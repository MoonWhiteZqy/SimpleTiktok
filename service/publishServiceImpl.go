package service

import "simpleTiktok/dao"

type PublishServiceImpl struct{}

func (p PublishServiceImpl) SavePOSTFile(content []byte, path, fileName, title string, userId int64) error {
	return dao.SavePOSTFile(content, path, fileName, title, userId)
}

func (p PublishServiceImpl) GetVideoOfUser(userId int64) ([]dao.Video, error) {
	return dao.GetVideoOfUser(userId)
}

func (p PublishServiceImpl) GetFeed() ([]dao.Video, error) {
	return dao.GetFeed()
}
