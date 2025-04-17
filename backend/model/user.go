package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"account-system/common"
)

// User 用戶模型
type User struct {
	Id               int            `json:"id"`
	Username         string         `json:"username" gorm:"unique;index" validate:"max=12"`
	Password         string         `json:"password" gorm:"not null;" validate:"min=8,max=20"`
	DisplayName      string         `json:"display_name" gorm:"index" validate:"max=20"`
	Role             int            `json:"role" gorm:"type:int;default:1"`   // admin, common
	Status           int            `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Email            string         `json:"email" gorm:"index" validate:"max=50"`
	AccessToken      *string        `json:"access_token" gorm:"type:char(32);column:access_token;uniqueIndex"` // 系統管理令牌
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Setting          string         `json:"setting" gorm:"type:text;column:setting"`
}

// UserBase 用戶基本信息，用於緩存
type UserBase struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Status   int    `json:"status"`
	Group    string `json:"group"`
	Quota    int    `json:"quota"`
	Setting  string `json:"setting"`
	Email    string `json:"email"`
}

func (user *User) ToBaseUser() *UserBase {
	cache := &UserBase{
		Id:       user.Id,
		Username: user.Username,
		Status:   user.Status,
		Email:    user.Email,
	}
	return cache
}

func (user *User) GetAccessToken() string {
	if user.AccessToken == nil {
		return ""
	}
	return *user.AccessToken
}

func (user *User) SetAccessToken(token string) {
	user.AccessToken = &token
}

// Insert 插入新用戶
func (user *User) Insert() error {
	var err error
	if user.Password != "" {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	user.SetAccessToken(common.GetUUID())
	result := DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update 更新用戶信息
func (user *User) Update(updatePassword bool) error {
	var err error
	if updatePassword {
		user.Password, err = common.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}

	newUser := *user
	updates := map[string]interface{}{
		"username":     newUser.Username,
		"display_name": newUser.DisplayName,
		"email":        newUser.Email,
	}
	if updatePassword {
		updates["password"] = newUser.Password
	}

	DB.First(&user, user.Id)
	if err = DB.Model(user).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// Delete 刪除用戶
func (user *User) Delete() error {
	if user.Id == 0 {
		return errors.New("id 為空！")
	}
	if err := DB.Delete(user).Error; err != nil {
		return err
	}
	return nil
}

// ValidateAndFill 驗證用戶密碼並填充用戶信息
func (user *User) ValidateAndFill() (err error) {
	password := user.Password
	username := strings.TrimSpace(user.Username)
	if username == "" || password == "" {
		return errors.New("用戶名或密碼為空")
	}
	// 通過用戶名或郵箱查找用戶
	DB.Where("username = ? OR email = ?", username, username).First(user)
	okay := common.ValidatePasswordAndHash(password, user.Password)
	if !okay || user.Status != common.UserStatusEnabled {
		return errors.New("用戶名或密碼錯誤，或用戶已被封禁")
	}
	return nil
}

// FillUserById 通過 ID 填充用戶信息
func (user *User) FillUserById() error {
	if user.Id == 0 {
		return errors.New("id 為空！")
	}
	err := DB.First(user, "id = ?", user.Id).Error
	return err
}

// GetUserById 通過 ID 獲取用戶
func GetUserById(id int, selectAll bool) (*User, error) {
	if id == 0 {
		return nil, errors.New("id 為空！")
	}
	user := User{Id: id}
	var err error = nil
	if selectAll {
		err = DB.First(&user, "id = ?", id).Error
	} else {
		err = DB.Omit("password").First(&user, "id = ?", id).Error
	}
	return &user, err
}

// DeleteUserById 通過 ID 刪除用戶
func DeleteUserById(id int) (err error) {
	if id == 0 {
		return errors.New("id 為空！")
	}
	user := User{Id: id}
	return user.Delete()
}

// GetAllUsers 獲取所有用戶
func GetAllUsers(page, pageSize int) (users []*User, total int64, err error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return nil, 0, err
	}

	// 獲取總數
	err = tx.Model(&User{}).Count(&total).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 獲取分頁數據
	err = tx.Omit("password").Order("id desc").Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 提交事務
	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// SearchUsers 搜索用戶
func SearchUsers(keyword string, page, pageSize int) (users []*User, total int64, err error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err = tx.Error; err != nil {
		return nil, 0, err
	}

	query := DB.Model(&User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR display_name LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 獲取總數
	err = query.Count(&total).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 獲取分頁數據
	err = query.Omit("password").Order("id desc").Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 提交事務
	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CheckUserExistOrDeleted 檢查用戶是否存在或已刪除
func CheckUserExistOrDeleted(username, email string) (bool, error) {
	if username == "" && email == "" {
		return false, errors.New("用戶名和郵箱均為空！")
	}
	var count int64
	var err error
	if username != "" && email != "" {
		err = DB.Unscoped().Model(&User{}).Where("username = ? OR email = ?", username, email).Count(&count).Error
	} else if username != "" {
		err = DB.Unscoped().Model(&User{}).Where("username = ?", username).Count(&count).Error
	} else {
		err = DB.Unscoped().Model(&User{}).Where("email = ?", email).Count(&count).Error
	}
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ValidateAccessToken 驗證訪問令牌
func ValidateAccessToken(token string) (user *User) {
	if token == "" {
		return nil
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	user = &User{}
	if DB.Where("access_token = ?", token).First(user).RowsAffected == 1 {
		return user
	}
	return nil
}

// IsUserEnabled 檢查用戶是否啟用
func IsUserEnabled(id int) (bool, error) {
	var user User
	err := DB.Where("id = ?", id).Select("status").Find(&user).Error
	if err != nil {
		return false, err
	}
	return user.Status == common.UserStatusEnabled, nil
}

// IsAdmin 檢查用戶是否為管理員
func IsAdmin(id int) bool {
	var user User
	err := DB.Where("id = ?", id).Select("role").Find(&user).Error
	if err != nil {
		return false
	}
	return user.Role >= common.RoleAdminUser
}

// IsRoot 檢查用戶是否為超級管理員
func IsRoot(id int) bool {
	var user User
	err := DB.Where("id = ?", id).Select("role").Find(&user).Error
	if err != nil {
		return false
	}
	return user.Role >= common.RoleRootUser
}

// GetMaxUserId 獲取最大用戶 ID
func GetMaxUserId() int {
	var user User
	err := DB.Order("id desc").First(&user).Error
	if err != nil {
		return 0
	}
	return user.Id
}

// createRootAccountIfNeed 如果需要，創建根用戶帳號
func createRootAccountIfNeed() error {
	var user User
	if err := DB.First(&user).Error; err != nil {
		common.SysLog("no user exists, create a root user for you: username is root, password is 123456")
		hashedPassword, err := common.Password2Hash("123456")
		if err != nil {
			return err
		}
		rootUser := User{
			Username:    "root",
			Password:    hashedPassword,
			Role:        common.RoleRootUser,
			Status:      common.UserStatusEnabled,
			DisplayName: "Root User",
			AccessToken: nil,
		}
		DB.Create(&rootUser)
	}
	return nil
}
