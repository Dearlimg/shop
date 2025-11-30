package dao

import (
	"database/sql"
	"shop/global/db"
	"shop/model"
)

// GetUserByUsername 根据用户名获取用户
func GetUserByUsername(username string) (*model.User, string, error) {
	var user model.User
	var passwordHash string
	err := db.DB.QueryRow(
		"SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &passwordHash, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

// CheckUsernameExists 检查用户名是否存在
func CheckUsernameExists(username string) (bool, error) {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return count > 0, nil
}

// CreateUser 创建用户
func CreateUser(username, password, email string) (int64, error) {
	result, err := db.DB.Exec(
		"INSERT INTO users (username, password, email) VALUES (?, ?, ?)",
		username, password, email,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
