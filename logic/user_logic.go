package logic

import (
	"fmt"

	"shop/dao"
	"shop/model"

	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func Register(req *model.RegisterRequest) (int64, error) {
	// 检查用户名是否已存在
	exists, err := dao.CheckUsernameExists(req.Username)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("密码加密失败: %w", err)
	}

	// 创建用户
	userID, err := dao.CreateUser(req.Username, string(hashedPassword), req.Email)
	if err != nil {
		return 0, fmt.Errorf("注册失败: %w", err)
	}

	return userID, nil
}

// Login 用户登录
func Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// 查询用户
	user, passwordHash, err := dao.GetUserByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("数据库查询错误: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 生成token（简化版，实际应使用JWT）
	token := fmt.Sprintf("user_%d", user.ID)

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
