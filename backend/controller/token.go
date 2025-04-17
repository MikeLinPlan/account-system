package controller

import (
	"account-system/common"
	"account-system/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetAllTokens 獲取所有令牌
func GetAllTokens(c *gin.Context) {
	userId := c.GetInt("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	tokens, total, err := model.GetTokensByUserId(userId, page, pageSize)
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
		"data":    tokens,
		"total":   total,
	})
}

// SearchTokens 搜索令牌
func SearchTokens(c *gin.Context) {
	userId := c.GetInt("id")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	tokens, total, err := model.SearchTokens(userId, keyword, page, pageSize)
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
		"data":    tokens,
		"total":   total,
	})
}

// GetToken 獲取特定令牌
func GetToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的令牌 ID",
		})
		return
	}
	token, err := model.GetTokenById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 檢查令牌是否屬於當前用戶
	userId := c.GetInt("id")
	if token.UserId != userId && !model.IsAdmin(userId) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無權訪問該令牌",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "獲取成功",
		"data":    token,
	})
}

// AddToken 添加令牌
func AddToken(c *gin.Context) {
	var token model.Token
	err := json.NewDecoder(c.Request.Body).Decode(&token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的參數",
		})
		return
	}
	userId := c.GetInt("id")
	token.UserId = userId
	token.Status = common.TokenStatusEnabled
	err = token.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "添加成功",
		"data":    token,
	})
}

// UpdateToken 更新令牌
func UpdateToken(c *gin.Context) {
	var token model.Token
	err := json.NewDecoder(c.Request.Body).Decode(&token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的參數",
		})
		return
	}
	if token.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "令牌 ID 不能為空",
		})
		return
	}
	// 檢查令牌是否存在
	existingToken, err := model.GetTokenById(token.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 檢查令牌是否屬於當前用戶
	userId := c.GetInt("id")
	if existingToken.UserId != userId && !model.IsAdmin(userId) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無權修改該令牌",
		})
		return
	}
	// 保留原始用戶 ID
	token.UserId = existingToken.UserId
	// 保留原始密鑰
	token.Key = existingToken.Key
	// 保留創建時間
	token.CreatedTime = existingToken.CreatedTime
	// 更新訪問時間
	token.AccessedTime = time.Now()
	err = token.Update()
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

// DeleteToken 刪除令牌
func DeleteToken(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無效的令牌 ID",
		})
		return
	}
	// 檢查令牌是否存在
	token, err := model.GetTokenById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	// 檢查令牌是否屬於當前用戶
	userId := c.GetInt("id")
	if token.UserId != userId && !model.IsAdmin(userId) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "無權刪除該令牌",
		})
		return
	}
	err = model.DeleteTokenById(id)
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
