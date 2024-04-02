package model

import (
	"encoding/json"
	"fmt"
	k8s "gin-template/model/kubernetes"
	"time"
)

type UserContainer struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	Label        string    `gorm:"size:255" json:"label"`
	Namespace    string    `gorm:"type:char(36)" json:"namespace"`
	ConfigID     int       `gorm:"type:char(36)" json:"configId"`
	UserID       int       `gorm:"type:char(36)" json:"userId"`
	ImageID      int       `gorm:"type:char(36)" json:"imageId"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`
	TotalRuntime int       `gorm:"type:int" json:"totalRuntime"`
	LastBoot     time.Time `gorm:"not null" json:"lastBoot"`
	StartCMD     string    `gorm:"size:255" json:"startCmd"`
	Status       string    `gorm:"size:255" json:"status"`
	Envs         string    `gorm:"type:json" json:"envs"`
	Ports        string    `gorm:"type:json" json:"ports"`
	Service      string    `gorm:"type:json" json:"service"`
}

type StorageInfo struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	ContainerID  int    `gorm:"int" json:"containerId"` // 应该改为pvcname,我先不处理 处理牛魔，没这个你怎么知道这是哪个容器的
	PVCName      string `gorm:"size:255" json:"pvcName"`
	StorageClass string `gorm:"size:255" json:"storageClass"`
	AccessMode   string `gorm:"size:255" json:"accessMode"`
	Type         string `gorm:"size:255" json:"type"`
	Size         string `gorm:"size:255" json:"size"` // 单位为G
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
	Ports           string
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

func GetConfigByID(id int) (*ContainerConfig, error) {
	var config ContainerConfig
	err := DB.First(&config, id).Error
	return &config, err
}

// 不知道返回啥还，反正报错肯定要返回,所以暂时就返回一个报错了
func CreateInstance(conf *InstanceConfig) error {
	podconf := instanceConfigToPodInfo(conf)
	return k8s.NewService(podconf)
}

// func TestInstance(pod k8s.Pod) error {
// 	if pod.Rescourses.Volumes == nil {
// 		pod.Rescourses.Volumes = make([]k8s.Storage, 0)
// 		pod.Rescourses.Volumes = append(pod.Rescourses.Volumes, k8s.Storage{
// 			PVCName:      fmt.Sprint("pvc-", pod.Name),
// 			RomLimit:     "15Gi",
// 			MountPath:    "/home/default",
// 			AccessMode:   "ReadWriteOnce",
// 			StorageClass: "nfs-storage",
// 		})
// 	}

// 	// fmt.Print(pod.Marshal()
// 	{
// 		// print pod to debug
// 		data, _ := json.Marshal(pod)
// 		fmt.Println(string(data))
// 	}

// 	return k8s.NewService(pod)
// }

func TestInstancev2(podconfig k8s.PodConfig) error {

	var pod k8s.Pod

	pod.Name = podconfig.Name
	pod.NameSpace = podconfig.NameSpace
	pod.Ports = podconfig.Resourses.Ports
	if pod.Rescourses.Volumes == nil {
		pod.Rescourses.Volumes = make([]k8s.Storage, 0)
		pod.Rescourses.Volumes = append(pod.Rescourses.Volumes, k8s.Storage{
			PVCName:      fmt.Sprint("pvc-", pod.Name),
			RomLimit:     podconfig.Resourses.DefaultVolumeSize,
			MountPath:    "/home/default",
			AccessMode:   "ReadWriteOnce",
			StorageClass: "nfs-storage",
		})
	}
	resouse, err := GetConfigByID(podconfig.Resourses.ConfigID)
	if err != nil {
		return err
	}
	imageurl, err := GetImageUrlByID(podconfig.ImgID)
	if err != nil {
		return err
	}
	pod.ImgUrl = imageurl

	pod.Rescourses.CPULimit = resouse.CpuConf
	pod.Rescourses.GPULimit = resouse.GpuConf
	pod.Rescourses.RamLimit = resouse.MemoryConf

	{
		// print pod to debug
		data, _ := json.Marshal(pod)
		fmt.Println(string(data))
	}

	return k8s.NewService(pod)
}

func SetUserContainerStatus(id int, status string) error {
	return DB.Model(&UserContainer{}).Where("id = ?", id).Update("status", status).Error
}

func SaveCreateConfig(podConfig k8s.PodConfig, userid int) (int, error) {
	// 创建UserContainer记录
	userContainer := UserContainer{
		Label:     podConfig.Name,
		Namespace: podConfig.NameSpace,
		ConfigID:  podConfig.Resourses.ConfigID, // 确保ConfigID为字符串
		UserID:    userid,                       // 假设这里是从上下文或其它地方获取的用户ID
		ImageID:   podConfig.ImgID,              // 假设这里是已知的或从PodConfig中提取的镜像ID
		CreatedAt: time.Now(),
		LastBoot:  time.Now(),
		StartCMD:  "",     // 根据需要设置
		Status:    "stop", // 假设新创建的容器初始状态为running
		Envs:      "{}",   // Envs, Ports, Service等字段根据PodConfig设置
		Ports:     "{}",
		Service:   "{}",
		// Envs, Ports, Service等字段根据PodConfig设置
	}

	// 保存UserContainer到数据库
	if err := DB.Create(&userContainer).Error; err != nil {
		return -1, fmt.Errorf("failed to create user container: %w", err)
	}

	// 创建StorageInfo记录
	storageInfo := StorageInfo{
		ContainerID:  userContainer.ID, // 使用UserContainer的ID
		PVCName:      fmt.Sprintf("pvc-%v", podConfig.Name),
		Size:         podConfig.Resourses.DefaultVolumeSize,
		Path:         "/home/default",
		AccessMode:   "ReadWriteOnce",
		StorageClass: "nfs-storage",
	}

	// 保存StorageInfo到数据库
	if err := DB.Create(&storageInfo).Error; err != nil {
		return -1, fmt.Errorf("failed to create storage info: %w", err)
	}

	return userContainer.ID, nil
}
