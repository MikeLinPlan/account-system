package controller

import (
	"account-system/common"
	"account-system/model"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// LoginRequest 登入請求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// setupLogin 設置登入會話
func setupLogin(user *model.User, c *gin.Context) {
	session := sessions.Default(c)
	session.Set("id", user.Id)
	session.Set("username", user.Username)
	session.Set("role", user.Role)
	session.Set("status", user.Status)
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "保存會話失敗: " + err.Error(),
			"success": false,
		})
		return
	}
	cleanUser := model.User{
		Id:          user.Id,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		Status:      user.Status,
		Email:       user.Email,
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"success": true,
		"data":    cleanUser,
	})
}

// Login 用戶登入
func Login(c *gin.Context) {
	if !common.PasswordLoginEnabled {
		c.JSON(http.StatusOK, gin.H{
			"message": "管理員關閉了密碼登入",
			"success": false,
		})
		return
	}
	var loginRequest LoginRequest
	err := json.NewDecoder(c.Request.Body).Decode(&loginRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "無效的參數",
			"success": false,
		})
		return
	}
	username := loginRequest.Username
	password := loginRequest.Password
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "無效的參數",
			"success": false,
		})
		return
	}
	user := model.User{
		Username: username,
		Password: password,
	}
	err = user.ValidateAndFill()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}
	setupLogin(&user, c)
}

// Register 用戶註冊
func Register(c *gin.Context) {
	if !common.RegisterEnabled {
		c.JSON(http.StatusOK, gin.H{
			"message": "管理員關閉了新用戶註冊",
			"success": false,
		})
		return
	}
	if !common.PasswordRegisterEnabled {
		c.JSON(http.StatusOK, gin.H{
			"message": "管理員關閉了通過密碼進行註冊，請使用第三方帳戶驗證的形式進行註冊",
			"success": false,
		})
		return
	}
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的參數",
		})
		return
	}
	// 驗證用戶輸入
	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用戶名和密碼不能為空",
		})
		return
	}
	if len(user.Password) < 8 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "密碼長度不得小於 8 位",
		})
		return
	}
	if common.EmailVerificationEnabled {
		if user.Email == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "管理員開啟了郵箱驗證，請輸入郵箱地址",
			})
			return
		}
		// 這裡可以添加郵箱驗證碼檢查邏輯
	}
	exist, err := model.CheckUserExistOrDeleted(user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "數據庫錯誤，請稍後重試",
		})
		common.SysError(fmt.Sprintf("CheckUserExistOrDeleted error: %v", err))
		return
	}
	if exist {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用戶名已存在，或已註銷",
		})
		return
	}
	cleanUser := model.User{
		Username:    user.Username,
		Password:    user.Password,
		DisplayName: user.Username,
		Email:       user.Email,
		Role:        common.RoleCommonUser,
		Status:      common.UserStatusEnabled,
	}
	if err := cleanUser.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "註冊成功",
	})
}

// Logout 用戶登出
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

// GetSelf 獲取當前用戶信息
func GetSelf(c *gin.Context) {
	id := c.GetInt("id")
	user, err := model.GetUserById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "獲取成功",
		"data":    user,
	})
}

// UpdateSelf 更新當前用戶信息
func UpdateSelf(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的參數",
		})
		return
	}
	updatePassword := user.Password != ""
	if updatePassword && len(user.Password) < 8 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "密碼長度不得小於 8 位",
		})
		return
	}
	user.Id = c.GetInt("id")
	err = user.Update(updatePassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}

// DeleteSelf 刪除當前用戶
func DeleteSelf(c *gin.Context) {
	id := c.GetInt("id")
	err := model.DeleteUserById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "刪除成功",
	})
}

// GenerateAccessToken 生成訪問令牌
func GenerateAccessToken(c *gin.Context) {
	id := c.GetInt("id")
	user, err := model.GetUserById(id, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	token := common.GetUUID()
	user.SetAccessToken(token)
	err = user.Update(false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "生成成功",
		"data":    token,
	})
}

// GetAllUsers 獲取所有用戶（管理員）
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	users, total, err := model.GetAllUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "獲取成功",
		"data":    users,
		"total":   total,
	})
}

// SearchUsers 搜索用戶（管理員）
func SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	users, total, err := model.SearchUsers(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "獲取成功",
		"data":    users,
		"total":   total,
	})
}

// GetUser 獲取特定用戶（管理員）
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的用戶 ID",
		})
		return
	}
	user, err := model.GetUserById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "獲取成功",
		"data":    user,
	})
}

// CreateUser 創建用戶（管理員）
func CreateUser(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	user.Username = strings.TrimSpace(user.Username)
	if err != nil || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的參數",
		})
		return
	}
	if len(user.Password) < 8 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "密碼長度不得小於 8 位",
		})
		return
	}
	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}
	myRole := c.GetInt("role")
	if user.Role >= myRole {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無法創建權限大於等於自己的用戶",
		})
		return
	}
	exist, err := model.CheckUserExistOrDeleted(user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "數據庫錯誤，請稍後重試",
		})
		return
	}
	if exist {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用戶名已存在，或已註銷",
		})
		return
	}
	if user.Status == 0 {
		user.Status = common.UserStatusEnabled
	}
	err = user.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "創建成功",
	})
}

// UpdateUser 更新用戶（管理員）
func UpdateUser(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的參數",
		})
		return
	}
	if user.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用戶 ID 不能為空",
		})
		return
	}
	updatePassword := user.Password != ""
	if updatePassword && len(user.Password) < 8 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "密碼長度不得小於 8 位",
		})
		return
	}
	myRole := c.GetInt("role")
	existingUser, err := model.GetUserById(user.Id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if existingUser.Role >= myRole {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無法修改權限大於等於自己的用戶",
		})
		return
	}
	if user.Role >= myRole {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無法將用戶權限設為大於等於自己的權限",
		})
		return
	}
	err = user.Update(updatePassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}

// DeleteUser 刪除用戶（管理員）
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的用戶 ID",
		})
		return
	}
	myRole := c.GetInt("role")
	existingUser, err := model.GetUserById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if existingUser.Role >= myRole {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無法刪除權限大於等於自己的用戶",
		})
		return
	}
	err = model.DeleteUserById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "刪除成功",
	})
}
