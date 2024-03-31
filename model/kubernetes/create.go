package kubernetes

import (
	"context"
	"fmt"
	"gin-template/common"
	"net"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func Int32Ptr(i int32) *int32 { return &i }

func GetStatefulSetMetadata(name string, nameapce string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: nameapce,
	}
}
func GetPodTemplateMetadata(label map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Labels: label,
	}
}
func GetPodSelector(label map[string]string) metav1.LabelSelector {
	return metav1.LabelSelector{
		MatchLabels: label,
	}
}

func GenerateLabel(name string) map[string]string {
	return map[string]string{
		"app": name,
	}
}

func GetContainer(pod Pod) apiv1.Container {
	var ports []apiv1.ContainerPort
	for _, port := range pod.Ports {
		ports = append(ports, apiv1.ContainerPort{
			ContainerPort: port,
		})
	}

	var mounts []apiv1.VolumeMount

	for _, volume := range pod.Rescourses.Volumes {
		mounts = append(mounts, apiv1.VolumeMount{
			Name:      volume.PVCName,
			MountPath: volume.MountPath,
		})
	}

	return apiv1.Container{
		Name:  pod.Name,
		Image: pod.ImgUrl,
		Ports: ports,
		Resources: apiv1.ResourceRequirements{
			Limits: apiv1.ResourceList{
				apiv1.ResourceMemory: resource.MustParse(pod.Rescourses.RamLimit),
			},
		},
		VolumeMounts: mounts,
	}
}

func GetPVC(pod Pod) []apiv1.PersistentVolumeClaim {
	var pvcs []apiv1.PersistentVolumeClaim
	// for _, volume := range pod.resource.volumes {
	for _, volume := range pod.Rescourses.Volumes {
		pvcs = append(pvcs, apiv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: volume.PVCName,
				Annotations: map[string]string{
					"volume.beta.kubernetes.io/storage-class": volume.StorageClass,
				},
			},
			Spec: apiv1.PersistentVolumeClaimSpec{
				AccessModes: []apiv1.PersistentVolumeAccessMode{
					apiv1.PersistentVolumeAccessMode(volume.AccessMode),
				},
				Resources: apiv1.VolumeResourceRequirements{
					Requests: apiv1.ResourceList{
						apiv1.ResourceStorage: resource.MustParse(volume.RomLimit),
					},
				},
			},
		})
	}
	return pvcs
}

func checkPort(port int32) bool {
	addr := ":" + strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func getAvailablePorts(l int32, r int32, cnt int) []int32 {
	ret := make([]int32, 0)
	ok := 1
	for i := l; i < r && ok < cnt; i++ {
		if checkPort(i) {
			ret = append(ret, i)
			ok++
		}
	}
	return ret
}

func GetServicePorts(pod Pod) []apiv1.ServicePort {
	forwardPorts := getAvailablePorts(30000, 40000, len(pod.Ports))
	ports := make([]apiv1.ServicePort, len(pod.Ports))
	for i, port := range pod.Ports {
		ports[i] = apiv1.ServicePort{
			Port:       int32(port),
			TargetPort: intstr.FromInt32(port),
			NodePort:   int32(forwardPorts[i]),
		}
	}
	return ports
}

func newStatefulSetAndService(podConf Pod) (*appsv1.StatefulSet, *apiv1.Service) {
	name := podConf.Name
	metadata := GetStatefulSetMetadata(name, podConf.NameSpace)
	label := GenerateLabel(name)

	selector := GetPodSelector(label)
	var containers []apiv1.Container
	containers = append(containers, GetContainer(podConf))
	ss := &appsv1.StatefulSet{
		ObjectMeta: metadata,
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    Int32Ptr(1),
			Selector:    &selector,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: GetPodTemplateMetadata(label),
				Spec: apiv1.PodSpec{
					Containers: containers,
				},
			},
			VolumeClaimTemplates: GetPVC(podConf),
		},
	}

	svc := &apiv1.Service{
		ObjectMeta: metadata,
		Spec: apiv1.ServiceSpec{
			Type:     apiv1.ServiceTypeNodePort,
			Ports:    GetServicePorts(podConf),
			Selector: label,
		},
	}

	return ss, svc
}

func NewService(pod Pod) error {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return err
	}

	namespace := pod.NameSpace
	statefulSetClient := clientset.AppsV1().StatefulSets(namespace) // 更新为 AppsV1
	serviceClient := clientset.CoreV1().Services(namespace)
	ss, svc := newStatefulSetAndService(pod)

	common.SysLog("Creating StatefulSet...")

	resultSts, err := statefulSetClient.Create(context.TODO(), ss, metav1.CreateOptions{})
	if err != nil {
		common.SysLog(err.Error())
		return err
		// panic(err)
	}
	{
		info := fmt.Sprintf("Created StatefulSet %q.\n", resultSts.GetObjectMeta().GetName())
		common.SysLog(info)
	}

	common.SysLog("Creating service...")

	resultSvc, err := serviceClient.Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		common.SysLog(err.Error())
		return err
	}
	{
		info := fmt.Sprintf("Created service %q.\n", resultSvc.GetObjectMeta().GetName())
		common.SysLog(info)
	}

	return nil
}

// func NewStatefulSetAndService(name string, img_url string, memReq string, memLim string, PVCMem string) (*appsv1.StatefulSet, *apiv1.Service) {
// 	ss := &appsv1.StatefulSet{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: name,
// 		},
// 		Spec: appsv1.StatefulSetSpec{
// 			ServiceName: name,
// 			Replicas:    Int32Ptr(1),
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: map[string]string{
// 					"app": "desktop",
// 				},
// 			},
// 			Template: apiv1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: map[string]string{
// 						"app": "desktop",
// 					},
// 				},
// 				Spec: apiv1.PodSpec{
// 					Containers: []apiv1.Container{
// 						{
// 							Name:  name,
// 							Image: img_url,
// 							Ports: []apiv1.ContainerPort{
// 								{
// 									ContainerPort: 3389,
// 								},
// 							},
// 							Resources: apiv1.ResourceRequirements{
// 								// Requests: apiv1.ResourceList{
// 								// 	apiv1.ResourceMemory: resource.MustParse(memReq),
// 								// },
// 								Limits: apiv1.ResourceList{
// 									apiv1.ResourceMemory: resource.MustParse(memLim),
// 								},
// 							},
// 							VolumeMounts: []apiv1.VolumeMount{
// 								{
// 									Name:      name + "-storage",
// 									MountPath: "/home/default",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			VolumeClaimTemplates: []apiv1.PersistentVolumeClaim{
// 				{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: name + "-storage",
// 						Annotations: map[string]string{
// 							"volume.beta.kubernetes.io/storage-class": "nfs-storage",
// 						},
// 					},
// 					Spec: apiv1.PersistentVolumeClaimSpec{
// 						AccessModes: []apiv1.PersistentVolumeAccessMode{
// 							apiv1.ReadWriteOnce,
// 						},
// 						Resources: apiv1.VolumeResourceRequirements{
// 							Requests: apiv1.ResourceList{
// 								apiv1.ResourceStorage: resource.MustParse(PVCMem),
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	svc := &apiv1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: name,
// 		},
// 		Spec: apiv1.ServiceSpec{
// 			Type: apiv1.ServiceTypeNodePort,
// 			Ports: []apiv1.ServicePort{
// 				{
// 					Port:       3389,
// 					TargetPort: intstr.FromInt(3389),
// 					NodePort:   30003,
// 				},
// 			},
// 			Selector: map[string]string{
// 				"app": name,
// 			},
// 		},
// 	}

// 	return ss, svc
// }
