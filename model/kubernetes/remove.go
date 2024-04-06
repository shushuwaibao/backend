package kubernetes

import (
	"context"
	"gin-template/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RemoveStatefulSet(name string, namespace string) error {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return err
	}
	_, err = clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err == nil {
		err := clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		return err
	} else {
		return err
	}
}

func RemovePVC(name string) error {
	
}
