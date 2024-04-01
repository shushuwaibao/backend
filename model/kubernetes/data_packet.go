package kubernetes

// Add missing imports here

type Storage struct {
	PVCName      string `json:"pvc_name"`
	RomLimit     string `json:"rom_limit" default:"15Gi"`
	MountPath    string `json:"mount_path" default:"/home/default"`
	AccessMode   string `json:"access_mode" default:"ReadWriteOnce"`
	StorageClass string `json:"storage_class" default:"nfs-storage"`
}

type Resource struct {
	Volumes  []Storage `json:"volumes"`
	CPULimit string    `json:"cpu_limit"`
	GPULimit string    `json:"gpu_limit"`
	RamLimit string    `json:"ram_limit" default:"2Gi"`
}

type Pod struct {
	Name       string   `json:"name" binding:"required"`
	NameSpace  string   `json:"namespace" default:"default"`
	ImgUrl     string   `json:"img_url" default:"172.16.13.73:18443/wb/lubuntu:v1.3"`
	Rescourses Resource `json:"resources"`
	Ports      []int32  `json:"ports" defalut:"{3398,22}"`
}
