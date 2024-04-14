package sftp

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SftpUpload(c *gin.Context) {
	var params FileTransferParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Printf("Bind json error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	header, err := c.FormFile("file")
	if err != nil {
		log.Printf("Get file header from context error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	file, err := header.Open()
	if err != nil {
		log.Printf("Open upload file error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open the file"})
		return
	}
	defer file.Close()

	err = uploadFile(&params, file)
	if err != nil {
		log.Printf("Upload file error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File uploaded successfully"})
}

func SftpDownload(c *gin.Context) {
	user := c.Query("user")
	password := c.Query("password")
	server := c.Query("server_with_ssh_port")
	path := c.Query("target_path")

	if user == "" || password == "" || server == "" || path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	params := &FileTransferParams{
		User: user,
		Password: password,
		ServerWithSshPort: server,
		TargetPath: path,
	}

	file, err := downloadFile(params)
	if err != nil {
		log.Printf("Download file error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file", "details": err.Error()})
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	if _, err = io.Copy(c.Writer, file); err != nil {
		log.Printf("Send downloaded file error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send the file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File downloaded successfully"})
}
