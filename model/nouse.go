package model

// type InstanceConfig struct {
// 	ID              int
// 	Label           string
// 	Namespace       string
// 	ConfigID        string
// 	UserID          string
// 	ImageID         string
// 	CreatedAt       time.Time
// 	TotalRuntime    int
// 	LastBoot        time.Time
// 	StartCMD        string
// 	Status          string
// 	Ports           string
// 	Envs            string
// 	Service         string
// 	AttachedStorage []StorageInfo // Assuming this field represents the joined data from another table.
// 	CPUConf         string        // Assuming you might also want the config details like CPU, Memory etc.
// 	MemoryConf      string
// 	GPUConf         string
// 	ImageConfig     ImageConfig // Assuming you want details about the image used.
// }

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

// func GetInstanceConfigByInstanceID(id int64) (*InstanceConfig, error) {
// 	var instanceConfig InstanceConfig
// 	err := DB.Table("user_containers").
// 		Select("user_containers.*, container_configs.cpu_conf, container_configs.memory_conf, container_configs.gpu_conf, image_configs.*").
// 		Joins("left join storage_infos on storage_infos.container_id = user_containers.id").
// 		Joins("left join container_configs on container_configs.config_id = user_containers.config_id").
// 		Joins("left join image_configs on image_configs.id = user_containers.image_id").
// 		Where("user_containers.id = ?", id).
// 		Scan(&instanceConfig).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Assuming multiple storages can be attached, find and assign them separately
// 	var attachedStorages []StorageInfo
// 	err = DB.Where("container_id = ?", id).Find(&attachedStorages).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	instanceConfig.AttachedStorage = attachedStorages

// 	return &instanceConfig, nil
// }