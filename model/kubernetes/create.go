package kubernetes

import (
	"context"
	"fmt"
	"gin-template/common"

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
			Requests: apiv1.ResourceList{
				apiv1.ResourceMemory: resource.MustParse("512Mi"),
			},
		},
		Env:          []apiv1.EnvVar{
			{
				Name: "INPUT_USER_NAME",
				Value: pod.Env.Uname,
			},
			{
				Name: "INPUT_USER_PSWD",
				Value: pod.Env.Password,
			}
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

func checkPort(targetPort int32) bool {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, service := range services.Items {
		// print(service.Name)
		common.SysLog(fmt.Sprintf("service name: %s\n", service.Name))
		for _, port := range service.Spec.Ports {
			if port.NodePort == targetPort {
				// fmt.Printf("Service %s in namespace %s is using port %d\n", service.Name, service.Namespace, checkPort)
				common.SysLog(fmt.Sprintf("Service %s in namespace %s is using port %d\n", service.Name, service.Namespace, targetPort))
				return false
			}
		}
	}
	return true
}

func getAvailablePorts(l int32, r int32, cnt int) []int32 {
	ret := make([]int32, 0)
	ok := 0
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
	common.SysLog(fmt.Sprintf("forwardPorts: %v\n", forwardPorts))
	ports := make([]apiv1.ServicePort, len(pod.Ports))
	for i, port := range pod.Ports {
		ports[i] = apiv1.ServicePort{
			Name:       fmt.Sprintf("from%vto%v", port, forwardPorts[i]),
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

	statefulSetClient := clientset.AppsV1().StatefulSets(namespace)
	ss, svc := newStatefulSetAndService(pod)

	common.SysLog("Checking if StatefulSet exists...")
	existSts, err := statefulSetClient.Get(context.TODO(), ss.Name, metav1.GetOptions{})

	if err == nil {
		common.SysLog(fmt.Sprintf("StatefulSet %s already exists, deleting...", existSts.Name))
		err := statefulSetClient.Delete(context.TODO(), ss.Name, metav1.DeleteOptions{})
		if err != nil {
			common.SysLog(err.Error())
			return err
		} else {
			info := fmt.Sprintf("Deleted StatefulSet %q.\n", existSts.GetObjectMeta().GetName())
			common.SysLog(info)
		}
	}

	{
		common.SysLog("Creating StatefulSet...")
		resultSts, err := statefulSetClient.Create(context.TODO(), ss, metav1.CreateOptions{})
		if err != nil {
			common.SysLog(err.Error())
			return err
		} else {
			info := fmt.Sprintf("Created StatefulSet %q.\n", resultSts.GetObjectMeta().GetName())
			common.SysLog(info)
		}
	}

	common.SysLog("Checking service...")
	serviceClient := clientset.CoreV1().Services(namespace)
	exsitSvc, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
	if err == nil {
		common.SysLog(fmt.Sprintf("service %s already exists, deleting ...", exsitSvc.Name))
		err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
		if err != nil {
			common.SysLog(err.Error())
			return err
		} else {
			info := fmt.Sprintf("Deleted service %q.\n", exsitSvc.GetObjectMeta().GetName())
			common.SysLog(info)
		}
	}

	{
		common.SysLog("creating service ...")
		resultSvc, err := serviceClient.Create(context.TODO(), svc, metav1.CreateOptions{})
		if err != nil {
			common.SysLog(err.Error())
			return err
		} else {
			info := fmt.Sprintf("Created service %q.\n", resultSvc.GetObjectMeta().GetName())
			common.SysLog(info)
		}
	}
	return nil
}
