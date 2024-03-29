package controller

import (
	"gin-template/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartInstanceByInstanceID(c *gin.Context) {

}

func StopInstanceByInstanceID(c *gin.Context) {

}

func RemoveInstancerByInstanceID(c *gin.Context) {

}

func CreateInstanceConfigAndStart(c *gin.Context) {

}

func EditInstanceConfig(c *gin.Context) {

}

func ExportInstanceImage(c *gin.Context) {

}

func GetAllAvailableInstanceConfig(c *gin.Context) {
	configs, err := model.GetAvailableInstanceConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, configs)
}
