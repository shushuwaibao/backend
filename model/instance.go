package model

import (
	"encoding/json"
	"fmt"
	"gin-template/common"
	k8s "gin-template/model/kubernetes"
	"time"
)

type UserContainer struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	Label        string    `gorm:"size:255 unique" json:"label"`
	Namespace    string    `gorm:"type:char(36)" json:"namespace"`
	ConfigID     int       `gorm:"type:char(36)" json:"configId"`
	UserID       int       `gorm:"type:char(36)" json:"userId"`
	ImageID      int       `gorm:"type:char(36)" json:"imageId"`
	CreatedAt    time.Time `gorm:"type:datetime(3)" json:"createdAt"`
	TotalRuntime uint64    `gorm:"type:int" json:"totalRuntime"`
	LastBoot     time.Time `gorm:"type:datetime(3)" json:"lastBoot"`
	StartCMD     string    `gorm:"size:255" json:"startCmd"`
	Status       string    `gorm:"size:255" json:"status"`
	ClusterIP    string    `gorm:"size:255" json:"clusterIP"`
	Envs         string    `gorm:"type:json" json:"envs"`
	Ports        string    `gorm:"type:json" json:"ports"`
	Service      string    `gorm:"size:255" json:"service"`
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

func GetUserContainerByID(id int) (*UserContainer, error) {
	var container UserContainer
	err := DB.First(&container, id).Error
	return &container, err
}

func GetUserContainerByUserID(userID int) ([]UserContainer, error) {
	var containers []UserContainer
	err := DB.Where("user_id = ?", userID).Find(&containers).Error
	return containers, err
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

func CreateService(podconfig k8s.PodConfig) error {
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
func updContainer(container *UserContainer, status string) error {
	if container.Status == "removed" {
		// removed
		return fmt.Errorf("container has been removed")
	} else if container.Status != status {
		if status == "running" {
			container.LastBoot = time.Now()
		} else {
			container.TotalRuntime += uint64(time.Now().Sub(container.LastBoot).Seconds())
			container.LastBoot = time.Now()
		}
		container.Status = status
		return nil
	} else {
		return nil
	}
}

func SetUserContainerStatus(id int, status string) error {
	container, err := GetUserContainerByID(id)
	if err != nil {
		return err
	}
	err = updContainer(container, status)
	if err != nil {
		return err
	}
	return DB.Save(container).Error
}

func CreateInstance(podConfig k8s.PodConfig, userid int) (int, error) {
	// check weahter the label is unique
	db := DB.Begin()
	var count int64
	db.Model(&UserContainer{}).Where("label = ?", podConfig.Name).Count(&count)
	if count > 0 {
		return -1, fmt.Errorf("label %s is not unique", podConfig.Name)
	}

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
	if err := db.Create(&userContainer).Error; err != nil {
		db.Rollback()
		return -1, fmt.Errorf("failed to create user container: %w", err)
	}

	// 创建StorageInfo记录
	storageInfo := StorageInfo{
		// ContainerID:  userContainer.ID, // 使用UserContainer的ID
		PVCName: fmt.Sprintf("pvc-%v", podConfig.Name),
		Size:    podConfig.Resourses.DefaultVolumeSize,
		// Path:         "/home/default",
		AccessMode:   "ReadWriteOnce",
		StorageClass: "nfs-storage",
	}

	// 保存StorageInfo到数据库
	if err := db.Create(&storageInfo).Error; err != nil {
		db.Rollback()
		return -1, fmt.Errorf("failed to create storage info: %w", err)
	}

	acl := PVCACL{
		StorageID:  storageInfo.ID,
		UserID:     userid,
		Permission: "admin",
	}
	if err := db.Create(&acl).Error; err != nil {
		db.Rollback()
		return -1, fmt.Errorf("failed to create pvc acl: %w", err)
	}

	binding := StorageContainerBind{
		StorageID:   storageInfo.ID,
		ContainerID: userContainer.ID,
		MountPath:   "/home/default",
	}

	if err := db.Create(&binding).Error; err != nil {
		db.Rollback()
		return -1, fmt.Errorf("failed to create storage container bind: %w", err)
	}
	if err := CreateService(podConfig); err != nil {
		db.Rollback()
		return -1, fmt.Errorf("failed to create service: %w", err)
	}

	if err := db.Commit().Error; err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userContainer.ID, nil
}

func GetInstanceName(uid int, iid int) ([]string, error) {
	common.SysLog(fmt.Sprintf("uid: %d, iid: %d", uid, iid))
	if GetRight(uid, iid) == 0 {
		return nil, fmt.Errorf("No rights")
	} else {
		var results []struct {
			Label     string
			Namespace string
		}
		if err := DB.Table("user_containers").Select("label", "namespace").Where("id = ?", iid).Find(&results).Error; err != nil {
			return nil, err
		} else {
			if len(results) == 0 {
				return nil, fmt.Errorf("No such instance")
			} else if len(results) > 1 {
				return nil, fmt.Errorf("More than one instance found")
			} else {
				return []string{results[0].Label, results[0].Namespace}, nil
			}
		}

	}
}

func FlushInstanceConfig(cid int) error {
	var container UserContainer
	DB.Table("user_containers").Where("id = ?", cid).First(&container)

	if container.Status == "removed" {
		return nil
		// removed
		// return fmt.Errorf("container has been removed")
	}

	pod, err := k8s.GetSS(container.Label, container.Namespace)
	if err != nil {
		updContainer(&container, "removed")
		DB.Save(&container)
		return err
	}

	// upd lifetime
	if pod.Status.ReadyReplicas == 0 {
		updContainer(&container, "stop")
	} else if pod.Status.ReadyReplicas == 1 {
		updContainer(&container, "running")
	}

	svc, err := k8s.GetService(container.Label, container.Namespace)
	if err != nil {
		updContainer(&container, "removed")
		DB.Save(&container)
		return err
	}
	container.ClusterIP = svc.Spec.ClusterIP
	container.Service = fmt.Sprintf("%v:%v", svc.Spec.ClusterIP, svc.Spec.Ports[0].NodePort)
	var ports []struct {
		TargetPort  int32 `json:"targetPort"`
		ForwardPort int32 `json:"forwardPort"`
	}

	for _, port := range svc.Spec.Ports {
		ports = append(ports, struct {
			TargetPort  int32 `json:"targetPort"`
			ForwardPort int32 `json:"forwardPort"`
		}{
			TargetPort:  port.TargetPort.IntVal,
			ForwardPort: port.NodePort,
		})
	}

	portBytes, err := json.Marshal(ports)
	// fmt.Printf("%s", portBytes)
	if err != nil {
		return err
	}
	container.Ports = string(fmt.Sprintf("%s", portBytes))

	return DB.Save(&container).Error
}
