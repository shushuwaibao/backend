package model

import (
	"fmt"
	k8s "gin-template/model/kubernetes"
	"time"
)

type UserContainer struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	Label        string    `gorm:"size:255" json:"label"`
	Namespace    string    `gorm:"type:char(36)" json:"namespace"`
	ConfigID     string    `gorm:"type:char(36)" json:"configId"`
	UserID       string    `gorm:"type:char(36)" json:"userId"`
	ImageID      string    `gorm:"type:char(36)" json:"imageId"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`
	TotalRuntime int       `gorm:"type:int" json:"totalRuntime"`
	LastBoot     time.Time `json:"lastBoot"`
	StartCMD     string    `gorm:"size:255" json:"startCmd"`
	Status       string    `gorm:"size:255" json:"status"`
	Envs         string    `gorm:"type:json" json:"envs"`
	Service      string    `gorm:"type:json" json:"service"`
}

type StorageInfo struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	ContainerID  int    `gorm:"int" json:"containerId"` // 应该改为pvcname,我先不处理
	StorageClass string `gorm:"size:255" json:"storageClass"`
	Type         string `gorm:"size:255" json:"type"`
	Size         int    `gorm:"type:int" json:"size"` // 单位为G
	Path         string `gorm:"size:255" json:"path"`
	NodeID       int    `gorm:"int" json:"nodeId"`
}

type ContainerConfig struct {
	ConfigID    int     `gorm:"primaryKey" json:"configId"`
	Name        string  `gorm:"size:255" json:"name"`
	CpuConf     string  `gorm:"size:255" json:"cpuConf"`
	GpuConf     string  `gorm:"size:255" json:"gpuConf"`
	MemoryConf  string  `gorm:"size:255" json:"memoryConf"` //?这是不是有点抽象，应该是什么2G啊啥的
	DefaultSize string  `gorm:"size:255" json:"defaultSize"`
	Price       float64 `gorm:"type:float" json:"price"`
}

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
	Envs            string
	Service         string
	AttachedStorage []StorageInfo // Assuming this field represents the joined data from another table.
	CPUConf         string        // Assuming you might also want the config details like CPU, Memory etc.
	MemoryConf      string
	GPUConf         string
	ImageConfig     ImageConfig // Assuming you want details about the image used.
}

func instanceConfigToPodInfo(instanceConfig *InstanceConfig) k8s.Pod {
	var ret k8s.Pod
	// ret.resource.cpu_core_limit = instanceConfig.CPUConf
	// ret.resource.ram_limit = instanceConfig.MemoryConf
	// ret.resource.gpu_core_limit = instanceConfig.GPUConf
	// ret.resource.volumes = make([]k8s.Storage, len(instanceConfig.AttachedStorage))
	// for i, storage := range instanceConfig.AttachedStorage {
	// 	ret.resource.volumes[i] = k8s.Storage{
	// 		pvc_name:      fmt.Sprintf("pvc-%d", storage.ID),
	// 		memory_limit:  fmt.Sprintf("%dGi", storage.Size),
	// 		mount_path:    storage.Path,
	// 		access_mode:   "ReadWriteOnce",
	// 		storage_class: storage.StorageClass,
	// 	}
	// }
	// ret.name = fmt.Sprint("rdp_desktop", instanceConfig.ID)
	// ret.img_url = GetContainerUrl(&instanceConfig.ImageConfig)
	ret.Rescourses.CPULimit = instanceConfig.CPUConf
	ret.Rescourses.RamLimit = instanceConfig.MemoryConf
	ret.Rescourses.GPULimit = instanceConfig.GPUConf
	ret.Name = fmt.Sprint("rdp_desktop", instanceConfig.ID)
	ret.ImgUrl = GetContainerUrl(&instanceConfig.ImageConfig)
	ret.Ports = []int32{3389, 22}
	ret.Rescourses.Volumes = make([]k8s.Storage, len(instanceConfig.AttachedStorage))
	for i, storage := range instanceConfig.AttachedStorage {
		ret.Rescourses.Volumes[i] = k8s.Storage{
			PVCName:      fmt.Sprintf("pvc-%d", storage.ID),
			RomLimit:     fmt.Sprintf("%dGi", storage.Size),
			MountPath:    storage.Path,
			AccessMode:   "ReadWriteOnce",
			StorageClass: storage.StorageClass,
		}
	}

	return ret
}

func GetUserContainerByID(id int) (*UserContainer, error) {
	var container UserContainer
	err := DB.First(&container, id).Error
	return &container, err
}

func GetALLUserContainerByUserID(userID string) ([]UserContainer, error) {
	var containers []UserContainer
	err := DB.Where("user_id = ?", userID).Find(&containers).Error
	return containers, err
}

func GetInstanceConfigByInstanceID(id int64) (*InstanceConfig, error) {
	var instanceConfig InstanceConfig
	err := DB.Table("user_containers").
		Select("user_containers.*, container_configs.cpu_conf, container_configs.memory_conf, container_configs.gpu_conf, image_configs.*").
		Joins("left join storage_infos on storage_infos.container_id = user_containers.id").
		Joins("left join container_configs on container_configs.config_id = user_containers.config_id").
		Joins("left join image_configs on image_configs.id = user_containers.image_id").
		Where("user_containers.id = ?", id).
		Scan(&instanceConfig).Error

	if err != nil {
		return nil, err
	}

	// Assuming multiple storages can be attached, find and assign them separately
	var attachedStorages []StorageInfo
	err = DB.Where("container_id = ?", id).Find(&attachedStorages).Error
	if err != nil {
		return nil, err
	}
	instanceConfig.AttachedStorage = attachedStorages

	return &instanceConfig, nil
}

func GetAvailableInstanceConfig() ([]ContainerConfig, error) {
	var configs []ContainerConfig
	err := DB.Find(&configs).Error
	return configs, err
}

//不知道返回啥还，反正报错肯定要返回,所以暂时就返回一个报错了
func CreateInstance(conf *InstanceConfig) error {
	podconf := instanceConfigToPodInfo(conf)
	return k8s.NewService(podconf)
}
