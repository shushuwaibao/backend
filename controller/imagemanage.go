package controller

import (
	"gin-template/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAvailableArchiveHandler(c *gin.Context) {
	imageConfigs, err := model.GetAvailableArchive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
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

	if result >= 1 {
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

	if result >= 1 {
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

func AddNewImageHandler(c *gin.Context) {
	var ii model.ImageConfig // 声明ImageConfig变量以接收请求体中的数据

	// 绑定请求体到ImageConfig结构体
	if err := c.ShouldBindJSON(&ii); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "解析失败",
		})
		return
	}

	// 调用业务逻辑函数添加新图片
	affectedRows, err := model.AddNewImage(&ii)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	if affectedRows >= 1 {
		// 如果添加成功，返回成功状态及受影响的行数
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "添加成功",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "添加未成功执行",
		})
	}
}
