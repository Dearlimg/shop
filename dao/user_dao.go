package dao

import (
	"errors"
	"gorm.io/gorm"
	"shop/global/db"
	"shop/model"
)

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*model.User, string, error) {
	var user model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", nil
		}
		return nil, "", err
	}
	return &user, user.Password, nil
}

// CheckUsernameExists 检查用户名是否存在
func CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := db.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// CreateUser 创建用户
func CreateUser(username, password, email string) (int64, error) {
	user := model.User{
		Username: username,
		Password: password,
		Email:    email,
	}
	err := db.DB.Create(&user).Error
	if err != nil {
		return 0, err
	}
	return int64(user.ID), nil
}
