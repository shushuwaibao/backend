package controller

import (
	"gin-template/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartInstanceByInstanceID(c *gin.Context) {
	// 传入一个id，根据id启动一个实例

}

func StopInstanceByInstanceID(c *gin.Context) {
	// 传入一个id，根据id停止一个实例

}

func RemoveInstancerByInstanceID(c *gin.Context) {
	// 传入一个id，根据id删除一个实例

}

func CreateInstanceConfigAndStart(c *gin.Context) {
	// 传入一个配置并启动
	// 传入的配置是一个json，包含了实例的配置信息，格式未定，按照需要加键值对

}

func EditInstanceConfig(c *gin.Context) {
	// 传入一个id和一个修改后的配置json，根据id修改配置
}

func ExportInstanceImage(c *gin.Context) {
	// 获取container的id->利用其他机子上的程序调用nerdctl commit并且push到仓库
}

func GetAllAvailableInstanceConfig(c *gin.Context) {
	configs, err := model.GetAvailableInstanceConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, configs)
}
