package model

type PVCACL struct {
	StorageID    int    `gorm:"type:int;primaryKey" json:"storageId"`
	SharedUserID int    `gorm:"type:int;primaryKey" json:"sharedUserId"`
	Permission   string `gorm:"size:255" json:"permission"` // read(only the rdp session), write(modify), admin(delete)
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
