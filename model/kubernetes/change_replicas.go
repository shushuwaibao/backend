package kubernetes

import (
	"context"
	"gin-template/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ChangeReplicas(namespace string, name string, replicas int) error {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return err
	}
	ss, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err == nil {
		*ss.Spec.Replicas = int32(replicas)
		_, err := clientset.AppsV1().StatefulSets(namespace).Update(context.TODO(), ss, metav1.UpdateOptions{})
		return err
	} else {
		return err
	}
}
