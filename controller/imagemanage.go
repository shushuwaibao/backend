package controller

import (
	"gin-template/model"
	"net/http"

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
			"message": "删除成功",
		})
	}
}

type UpdateImagePermissionRequest struct {
	ImageList []model.ImageConfig `json:"imageList"`
	NewValue  []string            `json:"newValue"`
}

func UpdateImagePermissionHandler(c *gin.Context) {
	// 解析请求体中的JSON数据到imageList和newValue
	var request UpdateImagePermissionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "解析失败",
		})
		return
	}

	//赋值
	newValue := request.NewValue
	imageList := request.ImageList

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
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "更新未成功执行",
		})
	}

}
