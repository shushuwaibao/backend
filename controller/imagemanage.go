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
