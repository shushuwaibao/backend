package model

import (
	"fmt"
	k8s "gin-template/model/kubernetes"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
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
	CPUConf         int           // Assuming you might also want the config details like CPU, Memory etc.
	MemoryConf      int
	GPUConf         int
	ImageConfig     ImageConfig // Assuming you want details about the image used.
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

func NewStatefulSetAndService(instanceConfig *InstanceConfig) (*appsv1.StatefulSet, *apiv1.Service) {
	name := fmt.Sprint("a", instanceConfig.ID)
	metadata := k8s.CreateStatefulSetMetadata(name, instanceConfig.Namespace)
	selector := CreateSelector()
	containers := createContainers(name, img_url, memReq, memLim)
	pvc := createPVC(name, PVCMem)
	servicePorts := createServicePorts()
	serviceSelector := createServiceSelector(name)

	// 使用上面的子函数来构建 StatefulSet 和 Service
	ss := &appsv1.StatefulSet{
		ObjectMeta: metadata,
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    k8s.Int32Ptr(1),
			Selector:    selector,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "desktop",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: containers,
				},
			},
			VolumeClaimTemplates: []apiv1.PersistentVolumeClaim{pvc},
		},
	}

	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeNodePort,
			Ports:    servicePorts,
			Selector: serviceSelector,
		},
	}

	return ss, svc
}

// func CreatePodFromInstanceConfig(instanceConfig *InstanceConfig) (*corev1.Pod, error) {
// 	// 创建Kubernetes客户端
// 	clientset, err := kubernetes.NewForConfig(Kube_Config)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
// 	}

// 	var envVars []corev1.EnvVar
// 	if instanceConfig.Envs != "" {
// 		err := json.Unmarshal([]byte(instanceConfig.Envs), &envVars)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal environment variables: %w", err)
// 		}
// 	}
// 	var volumeMounts []corev1.VolumeMount
// 	var volumes []corev1.Volume

// 	for _, storage := range instanceConfig.AttachedStorage {
// 		volumeName := fmt.Sprintf("storage-%d", storage.ID) // 使用存储的ID作为卷的名称
// 		volumeMounts = append(volumeMounts, corev1.VolumeMount{
// 			Name:      volumeName,
// 			MountPath: storage.Path, // 使用StorageInfo中定义的挂载路径
// 		})
// 		volumes = append(volumes, corev1.Volume{
// 			Name: volumeName,
// 			VolumeSource: corev1.VolumeSource{
// 				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
// 					ClaimName: fmt.Sprintf("pvc-%d", storage.ID), // 假设PVC的命名规则为"pvc-"加上存储的ID
// 				},
// 			},
// 		})
// 	}

// 	pod := &corev1.Pod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "a" + fmt.Sprint(instanceConfig.ID), // 假设pod名称为"a"加上instance ID
// 			Namespace: instanceConfig.Namespace,
// 			Labels:    map[string]string{"app": "instance"},
// 		},
// 		Spec: corev1.PodSpec{
// 			Containers: []corev1.Container{
// 				{
// 					Name:    instanceConfig.Label,
// 					Image:   instanceConfig.ImageConfig.Name, // 假设ImageConfig.Name包含了完整的镜像地址
// 					Command: []string{"/bin/sh", "-c", instanceConfig.StartCMD},
// 					Resources: corev1.ResourceRequirements{
// 						Requests: corev1.ResourceList{
// 							corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%dm", instanceConfig.CPUConf)),
// 							corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dMi", instanceConfig.MemoryConf)),
// 						},
// 						Limits: corev1.ResourceList{
// 							corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%dm", instanceConfig.CPUConf*2)), // 假设limit是request的两倍
// 							corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dMi", instanceConfig.MemoryConf*2)),
// 						},
// 					},
// 					Env:          envVars,
// 					VolumeMounts: volumeMounts,
// 				},
// 			},
// 			Volumes: volumes,
// 		},
// 	}
// 	// 创建Pod
// 	pod, err = clientset.CoreV1().Pods(instanceConfig.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create pod: %w", err)
// 	}

// 	return pod, nil
// }
