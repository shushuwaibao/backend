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
	userID, exists := c.Get("id")
	if !exists {
		// 如果不存在，可能是因为用户未认证
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// var instanceID int
	var instanceID struct {
		Iid int `json:"iid"`
	}
	err := c.ShouldBindJSON(&instanceID)
	if err != nil {
		//没有传入iid
		common.SysError(fmt.Sprintf("error: %v, val", err, instanceID.Iid))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	strs, err := model.GetInstanceName(userID.(int), instanceID.Iid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Get instance Error", "info": err.Error()})
		return
	}
	err = k8s.ChangeReplicas(strs[0], strs[1], 1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Start instance Error", "info": err.Error()})
		return
	}

	err = model.SetUserContainerStatus(instanceID.Iid, "running")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Save status Error", "info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": "successfully started instance"})
}

func StopInstanceByInstanceID(c *gin.Context) {
	// 传入一个id，根据id停止一个实例
	userID, exists := c.Get("id")
	if !exists {
		// 如果不存在，可能是因为用户未认证
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var instanceID struct {
		Iid int `json:"iid"`
	}
	err := c.ShouldBindJSON(&instanceID)
	if err != nil {
		common.SysError(fmt.Sprintf("error: %v, val", err, instanceID.Iid))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	strs, err := model.GetInstanceName(userID.(int), instanceID.Iid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Get instance Error", "info": err.Error()})
		return
	}
	err = k8s.ChangeReplicas(strs[0], strs[1], 0)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Stop instance Error", "info": err.Error()})
		return
	}

	err = model.SetUserContainerStatus(instanceID.Iid, "stop")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Save status Error", "info": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": "successfully started instance"})
}

func RemoveInstancerByInstanceID(c *gin.Context) {
	// 传入一个id，根据id删除一个实例

}

func CreateInstanceConfigAndStartv3(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	common.SysLog("conifg binded")

	podconf.Name = fmt.Sprint(podconf.Name, "-", userID)
	common.SysLog("config name: " + podconf.Name)
	cid, err := model.SaveCreateConfig(podconf, userID.(int))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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

func ListAllInstance(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		// 如果不存在，可能是因为用户未认证
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	instances, err := model.GetUserContainerByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, instances)
}

func GetAllAvailableInstanceConfig(c *gin.Context) {
	configs, err := model.GetAvailableInstanceConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, configs)
}
