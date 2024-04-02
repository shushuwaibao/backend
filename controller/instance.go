package controller

import (
	"fmt"
	"gin-template/common"
	"gin-template/model"
	k8s "gin-template/model/kubernetes"

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

// func CreateInstanceConfigAndStart(c *gin.Context) {
// 	// 传入一个配置并启动
// 	// 传入的配置是一个json，包含了实例的配置信息，格式未定，按照需要加键值对
// 	common.SysLog("receive a request creating instance")
// 	var podconf k8s.Pod
// 	err := c.ShouldBindJSON(&podconf)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	common.SysLog("conifg binded")

// 	err = model.TestInstance(podconf)
// 	if err == nil {
// 		c.JSON(http.StatusAccepted, gin.H{"info": "successfully created service"})
// 	} else {
// 		c.String(http.StatusBadRequest, err.Error())
// 	}
// }

func CreateInstanceConfigAndStartv3(c *gin.Context) {
	/*

		{
		  "name": "test-api",
		  "namespace": "default",
		  "img_id": 1,
		  "resources": {
				"config_id" : 1,
				"default_volume_size": "15Gi"
				"Ports": [3389:22222]
		  }
		}


	*/
	userID, exists := c.Get("id")
	if !exists {
		// 如果不存在，可能是因为用户未认证
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// save config
	var podconf k8s.PodConfig
	err := c.ShouldBindJSON(&podconf)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	common.SysLog("conifg binded")

	podconf.Name = fmt.Sprint(podconf.Name, "-", userID)
	common.SysLog("config name: " + podconf.Name)
	cid, err := model.SaveCreateConfig(podconf, userID.(int))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = model.TestInstancev2(podconf)
	if err == nil {
		err = model.SetUserContainerStatus(cid, "running")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "config successfully saved but service failed"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"info": "successfully created service"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config successfully saved but service failed"})
	}
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
