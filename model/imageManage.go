package model

import (
	"errors"
	"log"
)

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
