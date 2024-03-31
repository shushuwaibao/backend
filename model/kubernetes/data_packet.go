package kubernetes

// Add missing imports here

type Storage struct {
	PVCName      string
	RomLimit     string
	MountPath    string
	AccessMode   string
	StorageClass string
}

type Resource struct {
	Volumes  []Storage
	CPULimit string
	GPULimit string
	RamLimit string
}

type Pod struct {
	Name       string
	ImgUrl     string
	Rescourses Resource
	Ports      []int32
}
