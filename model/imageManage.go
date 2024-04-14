package model

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type tag struct {
	Repository_id int    `json:"repository_id"`
	Name          string `json:"name"`
	Push_time     string `json:"push_time"`
	Pull_time     string `json:"pull_time"`
	Signed        bool   `json:"signed"`
	Id            int    `json:"id"`
	Immutable     bool   `json:"immutable"`
	Artifact_id   int    `json:"artifact_id"`
}

func GetAvailableArchive() ([]ImageConfig, error) {
	var results []ImageConfig
	// 使用已经初始化的 DB 进行查询
	err := DB.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func DeleteImage(imageList []ImageConfig) (int64, error) {
	// 遍历imageList，并删除每一个ImageConfig记录
	var deletedCount int64
	for _, image := range imageList {
		result := DB.Delete(&image, "id = ?", image.ID)
		if result.Error != nil {
			return deletedCount, result.Error
		}
		deletedCount += result.RowsAffected
	}
	return deletedCount, nil
}

func UpdateImagePermission(imageList []ImageConfig, newValue []string) (int64, error) {
	// 确保传入的列表长度一致
	if len(imageList) != len(newValue) {
		return 0, errors.New("the length of imageList and newValue must be the same")
	}

	// 遍历列表，更新每条记录
	var updatedRows int64
	for i := 0; i < len(imageList); i++ {
		// 使用gorm的Model结构中的ID字段作为条件更新permission字段
		result := DB.Model(&ImageConfig{}).Where("id = ?", imageList[i].ID).Update("permission", newValue[i])
		if result.Error != nil {
			log.Println("Error updating image permission:", result.Error)
			return updatedRows, result.Error
		}
		updatedRows += result.RowsAffected
	}

	return updatedRows, nil
}

func SearchImageIsNotExit(project_name, repository_name, version string) (int, error) {
	// 设置Harbor API的URL和认证信息
	harborURL := "https://172.16.13.73:18443/api/v2.0"
	username := "admin"
	password := "Harbor12345"
	var err error

	// 创建自定义的http.Transport，禁用TLS证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 创建HTTP客户端
	client := &http.Client{Transport: tr}

	// 创建请求
	req, err := http.NewRequest("GET", harborURL+"/projects/"+project_name+"/repositories/"+repository_name+"/artifacts/"+version+"/tags", nil) /**/

	if err != nil {
		return -1, err
	}

	// 设置基本认证头
	req.SetBasicAuth(username, password)

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return -2, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -3, err
	}

	var taglist []tag
	err = json.Unmarshal(body, &taglist)
	if err != nil {
		return -4, err
	}

	if len(taglist) == 1 {
		return 1, err
	}

	return -5, err
}

func AddNewImage(ii *ImageConfig) (int64, error) {
	// 从注册表的路径中提取项目名称
	project_name := strings.Split(ii.Registry, "/")[1]
	repository_name := ii.Name
	version := ii.Version

	// 检查图片是否已经存在
	size, err := SearchImageIsNotExit(project_name, repository_name, version)
	if err != nil {
		// 如果在检查过程中发生错误，则直接返回错误
		return 0, err
	}

	if size <= 0 {
		// 如果图片已经存在，则不需要进行任何操作，返回0和nil错误
		return 0, nil
	}
	// 如果size是正数，表示图片不存在，可以进行添加操作

	//使用gorm的DB实例创建新的ImageConfig记录
	result := DB.Create(ii)
	if result.Error != nil {
		// 如果在创建过程中发生错误，则返回错误
		return 0, result.Error
	}
	// 返回添加的记录数，通常应为1，除非有触发器或其他数据库逻辑改变了这一点
	return result.RowsAffected, nil
}
