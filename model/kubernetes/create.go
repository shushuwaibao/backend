package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Int32Ptr(i int32) *int32 { return &i }


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
// 								Requests: apiv1.ResourceList{
// 									apiv1.ResourceMemory: resource.MustParse(memReq),
// 								},
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

func CreateStatefulSetMetadata(name string, nameapce string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: nameapce,
	}
}
