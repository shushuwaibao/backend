package model

type UserPVCShare struct {
	StorageID    int    `gorm:"type:char(36);primaryKey" json:"storageId"`
	SharedUserID int    `gorm:"type:char(36);primaryKey" json:"sharedUserId"`
	Permission   string `gorm:"size:255" json:"permission"` // read(only the rdp session), write(modify), admin(delete)
}

type StorageContainerBind struct {
	StorageID   int    `gorm:"type:char(36);primaryKey" json:"storageId"`
	ContainerID int    `gorm:"type:char(36);primaryKey" json:"containerId"`
	MountPath   string `gorm:"size:255" json:"mountPath"`
}

type StorageInfo struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	PVCName      string `gorm:"size:255" json:"pvcName"`
	StorageClass string `gorm:"size:255" json:"storageClass"`
	AccessMode   string `gorm:"size:255" json:"accessMode"`
	Type         string `gorm:"size:255" json:"type"`
	Size         string `gorm:"size:255" json:"size"` // 单位为G
	NodeID       int    `gorm:"int" json:"nodeId"`
}
