package service

import "simpleTiktok/dao"

type RelationService interface {
	Action(int64, string, string) error
	MasterList(string) ([]dao.User, error)
	FollowerList(string) ([]dao.User, error)
}
