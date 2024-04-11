package controller

import (
	"gin-template/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAvailableArchiveHandler(c *gin.Context) {
	imageConfigs, err := model.GetAvailableArchive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, imageConfigs)
}

func DeleteImageHandler(c *gin.Context) {
	var imageList []model.ImageConfig
	if err := c.ShouldBindJSON(&imageList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "解析失败",
		})
		return
	}

	// 调用SetImagePermission函数
	result, err := model.DeleteImage(imageList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除失败",
		})
		return
	}

	if result == 1 {
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "更新成功",
		})
	}
}

func UpdateImagePermissionHandler(c *gin.Context) {
	// 解析请求体中的JSON数据到imageList和newValue
	var imageList []model.ImageConfig
	var newValue []string
	if err := c.ShouldBindJSON(&imageList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "解析失败",
		})
		return
	}

	// 假设newValue是以逗号分隔的字符串，需要将其拆分为切片
	newValueStr := c.Query("newValue") // 假设newValue通过查询参数传递，格式为逗号分隔的字符串
	newValue = strings.Split(newValueStr, ",")

	// 确保newValue和imageList长度一致
	if len(newValue) != len(imageList) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数不一致",
		})
		return
	}

	// 调用SetImagePermission函数
	result, err := model.UpdateImagePermission(imageList, newValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新失败",
		})
		return
	}

	if result == 1 {
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "更新成功",
		})
	}

}
