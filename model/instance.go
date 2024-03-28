package model

import "time"

type UserContainer struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
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
	ContainerID  int    `gorm:"int" json:"containerId"`
	StorageClass string `gorm:"size:255" json:"storageClass"`
	Type         string `gorm:"size:255" json:"type"`
	Size         int    `gorm:"type:int" json:"size"`
	Path         string `gorm:"size:255" json:"path"`
	NodeID       int    `gorm:"int" json:"nodeId"`
}

type ContainerConfig struct {
	ConfigID    int     `gorm:"primaryKey" json:"configId"`
	Name        string  `gorm:"size:255" json:"name"`
	CpuConf     int     `gorm:"type:int" json:"cpuConf"`
	GpuConf     int     `gorm:"type:int" json:"gpuConf"`
	MemoryConf  int     `gorm:"type:int" json:"memoryConf"`
	DefaultSize int     `gorm:"type:int" json:"defaultSize"`
	Price       float64 `gorm:"type:float" json:"price"`
}

func GetUserContainerByID(id int64) (*UserContainer, error) {
	var container UserContainer
	err := DB.First(&container, id).Error
	return &container, err
}

func GetALLUserContainerByUSERID(userID string) ([]UserContainer, error) {
	var containers []UserContainer
	err := DB.Where("user_id = ?", userID).Find(&containers).Error
	return containers, err
}

func GetInstanceConfigByContainerID(id int) (*ContainerConfig, error) {
	var config ContainerConfig
	// get config id
	var container UserContainer
	err := DB.First(&container, id).Error
	if err != nil {
		return nil, err
	}
	err = DB.First(&config, container.ConfigID).Error
	return &config, err
}

func GetAvailableInstanceConfig() ([]ContainerConfig, error) {
	var configs []ContainerConfig
	err := DB.Find(&configs).Error
	return configs, err
}
