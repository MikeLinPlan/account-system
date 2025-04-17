package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"account-system/common"
)

// Token 令牌模型
type Token struct {
	Id               int            `json:"id"`
	UserId           int            `json:"user_id" gorm:"index"`
	Key              string         `json:"key" gorm:"type:varchar(64);uniqueIndex"`
	Name             string         `json:"name" gorm:"type:varchar(64)"`
	Status           int            `json:"status" gorm:"type:int;default:1"`
	CreatedTime      time.Time      `json:"created_time" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	AccessedTime     time.Time      `json:"accessed_time" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	ExpiredTime      time.Time      `json:"expired_time" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	RemainQuota      int            `json:"remain_quota" gorm:"type:int;default:0"`
	UnlimitedQuota   bool           `json:"unlimited_quota" gorm:"type:tinyint(1);default:0"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

// Insert 插入新令牌
func (token *Token) Insert() error {
	token.Key = common.GetUUID()
	token.CreatedTime = time.Now()
	token.AccessedTime = time.Now()
	token.ExpiredTime = time.Now().AddDate(10, 0, 0) // 默認10年有效期
	result := DB.Create(token)
	return result.Error
}

// Update 更新令牌
func (token *Token) Update() error {
	result := DB.Model(token).Updates(map[string]interface{}{
		"name":            token.Name,
		"status":          token.Status,
		"expired_time":    token.ExpiredTime,
		"remain_quota":    token.RemainQuota,
		"unlimited_quota": token.UnlimitedQuota,
	})
	return result.Error
}

// Delete 刪除令牌
func (token *Token) Delete() error {
	if token.Id == 0 {
		return errors.New("id 為空！")
	}
	result := DB.Delete(token)
	return result.Error
}

// ValidateUserToken 驗證用戶令牌
func ValidateUserToken(key string) (*Token, error) {
	if key == "" {
		return nil, errors.New("令牌為空")
	}
	var token Token
	err := DB.Where("key = ?", key).First(&token).Error
	if err != nil {
		return nil, errors.New("無效的令牌")
	}
	if token.Status != common.TokenStatusEnabled {
		return &token, errors.New("令牌已被禁用")
	}
	if token.ExpiredTime.Before(time.Now()) {
		token.Status = common.TokenStatusExpired
		DB.Save(&token)
		return &token, errors.New("令牌已過期")
	}
	if !token.UnlimitedQuota && token.RemainQuota <= 0 {
		token.Status = common.TokenStatusExhausted
		DB.Save(&token)
		return &token, errors.New("令牌額度已用盡")
	}
	token.AccessedTime = time.Now()
	DB.Save(&token)
	return &token, nil
}

// GetTokenById 通過 ID 獲取令牌
func GetTokenById(id int) (*Token, error) {
	if id == 0 {
		return nil, errors.New("id 為空！")
	}
	token := Token{Id: id}
	err := DB.First(&token, "id = ?", id).Error
	return &token, err
}

// GetTokensByUserId 獲取用戶的所有令牌
func GetTokensByUserId(userId int, page, pageSize int) (tokens []*Token, total int64, err error) {
	if userId == 0 {
		return nil, 0, errors.New("用戶 ID 為空！")
	}
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
	err = tx.Model(&Token{}).Where("user_id = ?", userId).Count(&total).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 獲取分頁數據
	err = tx.Where("user_id = ?", userId).Order("id desc").Limit(pageSize).Offset(offset).Find(&tokens).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 提交事務
	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

// SearchTokens 搜索令牌
func SearchTokens(userId int, keyword string, page, pageSize int) (tokens []*Token, total int64, err error) {
	if userId == 0 {
		return nil, 0, errors.New("用戶 ID 為空！")
	}
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

	query := DB.Model(&Token{}).Where("user_id = ?", userId)
	if keyword != "" {
		query = query.Where("name LIKE ? OR key LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 獲取總數
	err = query.Count(&total).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 獲取分頁數據
	err = query.Order("id desc").Limit(pageSize).Offset(offset).Find(&tokens).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	// 提交事務
	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

// DeleteTokenById 通過 ID 刪除令牌
func DeleteTokenById(id int) error {
	if id == 0 {
		return errors.New("id 為空！")
	}
	token := Token{Id: id}
	return token.Delete()
}
