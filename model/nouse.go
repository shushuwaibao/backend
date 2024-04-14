package model

import (
	"strconv"
	"time"
)

type InstanceConfig struct {
	ID              int
	Label           string
	Namespace       string
	ConfigID        string
	UserID          string
	ImageID         string
	CreatedAt       time.Time
	TotalRuntime    int
	LastBoot        time.Time
	StartCMD        string
	Status          string
	Ports           string
	Envs            string
	Service         string
	AttachedStorage []StorageInfo // Assuming this field represents the joined data from another table.
	CPUConf         string        // Assuming you might also want the config details like CPU, Memory etc.
	MemoryConf      string
	GPUConf         string
	ImageConfig     ImageConfig // Assuming you want details about the image used.
}

// func instanceConfigToPodInfo(instanceConfig *InstanceConfig) k8s.Pod {
// 	var ret k8s.Pod
// 	ret.Rescourses.CPULimit = instanceConfig.CPUConf
// 	ret.Rescourses.RamLimit = instanceConfig.MemoryConf
// 	ret.Rescourses.GPULimit = instanceConfig.GPUConf
// 	ret.Name = fmt.Sprint("rdp_desktop", instanceConfig.ID)
// 	ret.ImgUrl = GetContainerUrl(&instanceConfig.ImageConfig)
// 	ret.Ports = []int32{3389, 22}
// 	ret.Rescourses.Volumes = make([]k8s.Storage, len(instanceConfig.AttachedStorage))
// 	for i, storage := range instanceConfig.AttachedStorage {
// 		ret.Rescourses.Volumes[i] = k8s.Storage{
// 			PVCName:  fmt.Sprintf("pvc-%d", storage.ID),
// 			RomLimit: fmt.Sprintf("%dGi", storage.Size),
// 			// MountPath:    storage.Path,
// 			AccessMode:   "ReadWriteOnce",
// 			StorageClass: storage.StorageClass,
// 		}
// 	}

// 	return ret
// }

func GetInstanceConfigByInstanceID(iid int) (*InstanceConfig, error) {
	var userContainer UserContainer
	err := DB.Table("user_containers").Where("ID = ?", iid).First(&userContainer).Error
	if err != nil {
		return nil, err
	}

	// 获取配置详情，使用正确的表名 container_configs
	var containerConfig ContainerConfig
	err = DB.Table("container_configs").Where("config_id = ?", userContainer.ConfigID).First(&containerConfig).Error
	if err != nil {
		return nil, err
	}

	// 获取存储信息，确保正确使用表名 storage_infos 和 storage_container_binds
	var storageInfos []StorageInfo
	err = DB.Table("storage_infos").
		Joins("JOIN storage_container_binds ON storage_infos.id = storage_container_binds.storage_id").
		Where("storage_container_binds.container_id = ?", iid).
		Find(&storageInfos).Error
	if err != nil {
		return nil, err
	}

	// 获取图像配置，使用正确的表名 image_configs
	var imageConfig ImageConfig
	err = DB.Table("image_configs").Where("id = ?", userContainer.ImageID).First(&imageConfig).Error
	if err != nil {
		return nil, err
	}
	// Construct InstanceConfig.
	instanceConfig := InstanceConfig{
		ID:              userContainer.ID,
		Label:           userContainer.Label,
		Namespace:       userContainer.Namespace,
		ConfigID:        strconv.Itoa(userContainer.ConfigID),
		UserID:          strconv.Itoa(userContainer.UserID),
		ImageID:         strconv.Itoa(userContainer.ImageID),
		CreatedAt:       userContainer.CreatedAt,
		TotalRuntime:    int(userContainer.TotalRuntime),
		LastBoot:        userContainer.LastBoot,
		StartCMD:        userContainer.StartCMD,
		Status:          userContainer.Status,
		Ports:           userContainer.Ports,
		Envs:            userContainer.Envs,
		Service:         userContainer.Service,
		AttachedStorage: storageInfos,
		CPUConf:         containerConfig.CpuConf,
		MemoryConf:      containerConfig.MemoryConf,
		GPUConf:         containerConfig.GpuConf,
		ImageConfig:     imageConfig,
	}
	return &instanceConfig, nil
}
