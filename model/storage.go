package model

import (
	"fmt"
	k8s "gin-template/model/kubernetes"
)

type PVCACL struct {
	StorageID  int    `gorm:"type:int;primaryKey" json:"storageId"`
	UserID     int    `gorm:"type:int;primaryKey" json:"UserId"`
	Permission string `gorm:"size:255" json:"permission"` // read(only the rdp session), write(modify), admin(delete)
}

type StorageContainerBind struct {
	StorageID   int    `gorm:"type:int;primaryKey" json:"storageId"`
	ContainerID int    `gorm:"type:int;primaryKey" json:"containerId"`
	MountPath   string `gorm:"size:255" json:"mountPath"`
}

type StorageInfo struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	PVCName      string `gorm:"size:255" json:"pvcName"`
	StorageClass string `gorm:"size:255" json:"storageClass"`
	AccessMode   string `gorm:"size:255" json:"accessMode"`
	Type         string `gorm:"size:255" json:"type"` // 存储类型，例如nfs
	Size         string `gorm:"size:255" json:"size"` // 存储空间大小，如4Gi
	NodeID       int    `gorm:"int" json:"nodeId"`
}

func GetStorageInfo(storageID int) (StorageInfo, error) {
	var storageInfo StorageInfo
	result := DB.Table("storage_info").First(&storageInfo, storageID)
	return storageInfo, result.Error
}

func ListBindedStorage(containerID int) ([]int, error) {
	var storageContainerBinds []StorageContainerBind
	result := DB.Table("storage_container_bindss").Where("container_id = ?", containerID).Find(&storageContainerBinds)
	var storageIDs []int
	for _, bind := range storageContainerBinds {
		storageIDs = append(storageIDs, bind.StorageID)
	}
	return storageIDs, result.Error
}

func ListAdminPVC(UserID int) ([]int, error) {
	var pvcACLs []PVCACL
	result := DB.Table("pvc_acl").Where("shared_user_id = ? AND permission = ?", UserID, "admin").Find(&pvcACLs)
	// result := DB.Where("shared_user_id = ? AND permission = ?", UserID, "admin").Find(&pvcACLs)
	var storageIDs []int
	for _, pvcACL := range pvcACLs {
		storageIDs = append(storageIDs, pvcACL.StorageID)
	}
	return storageIDs, result.Error
}

func ListOnlyBindedPVC(userID int, containerID int) ([]StorageInfo, error) {
	var storageContainerBinds []StorageContainerBind
	result := DB.Table("storage_container_binds").Where("container_id = ?", containerID).Find(&storageContainerBinds)
	if result.Error != nil {
		return nil, result.Error
	}
	admin, err := ListAdminPVC(userID)
	if err != nil {
		return nil, err
	}
	binds, err := ListBindedStorage(containerID)
	if err != nil {
		return nil, err
	}
	// 取admin 和binds 的交集
	vis := make(map[int]bool)
	for _, v := range admin {
		vis[v] = true
	}
	var rem []int
	for _, v := range binds {
		if vis[v] {
			rem = append(rem, v)
		}
	}

	var res []StorageInfo
	for _, bind := range rem {
		//检查在db中是否只有一条bind
		var tmp []StorageInfo
		err = DB.Table("storage_container_binds").Where("container_id = ? AND storage_id = ?", containerID, bind).Find(&tmp).Error
		if err != nil {
			return nil, err
		}
		if len(tmp) == 1 {
			storageInfo, err := GetStorageInfo(bind)
			if err != nil {
				return nil, err
			}
			res = append(res, storageInfo)
		}
	}
	return res, nil
}

func DeleteStorageEntries(storage StorageInfo) error {
	tx := DB.Begin()
	if err := tx.Table("pvc_acl").Where("storage_id = ?", storage.ID).Delete(&PVCACL{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete from PVCACL: %w", err)
	}

	if err := tx.Table("storage_container_binds").Where("storage_id = ?", storage.ID).Delete(&StorageContainerBind{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete from StorageContainerBind: %w", err)
	}

	if err := k8s.RemovePVC(storage.PVCName); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove PVC: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
