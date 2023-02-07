package dao

import (
	"errors"

	"gorm.io/gorm"
)

type UserModel struct {
	*gorm.Model
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

// 查看用户名是否已经被使用
func GetUsersByUserName(userName string) []UserModel {
	var users []UserModel
	DB.Where("name = ?", userName).Find(&users)
	return users
}

// 创建用户
func CreateUser(userName, password string) (int64, error) {
	sameNameUser := GetUsersByUserName(userName)
	if len(sameNameUser) > 0 {
		return -1, errors.New("username already in use")
	}
	curUser := UserModel{Name: userName, Password: password}
	DB.Create(&curUser)
	if DB.Error != nil {
		return -1, DB.Error
	}
	return int64(curUser.ID), nil
}

// 检查密码是否相同
func CheckPassword(userName, password string) (int64, error) {
	sameNameUser := GetUsersByUserName(userName)
	if len(sameNameUser) == 0 {
		return -1, errors.New("user not exist")
	}
	pwd := sameNameUser[0].Password
	if pwd != password {
		return -1, errors.New("password no equal")
	}
	return int64(sameNameUser[0].ID), nil
}

// 根据用户id,获取用户名信息
func GetUserByUserId(userId int64) (string, error) {
	var user UserModel
	DB.Where("id = ?", userId).Find(&user)
	if DB.Error != nil {
		return "", DB.Error
	}
	if user.Model == nil {
		return "", errors.New("user id not exists")
	}
	return user.Name, nil
}

// 根据用户id返回用户信息
func getUserById(userId int64) (User, error) {
	var user UserModel
	DB.Where("id = ?", userId).Find(&user)
	if DB.Error != nil {
		return User{}, DB.Error
	}
	if user.Model == nil {
		return User{}, errors.New("user id not exists")
	}
	return User{UserId: userId, Name: user.Name}, nil
}
