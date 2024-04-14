package kubernetes

import (
	"context"
	"gin-template/common"

	apiv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	scv1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func GetSS(name string, namespace string) (*apiv1.StatefulSet, error) {
	// common.SysLog("get ss, name: " + name + ", namespace: " + namespace)
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return nil, err
	}
	ss, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		common.SysLog(err.Error())
		return nil, err
	} else {
		return ss, err
	}
}

func GetService(name string, namespace string) (*corev1.Service, error) {
	// common.SysLog("get service")
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return nil, err
	}
	svc, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		common.SysLog(err.Error())
		return nil, err
	} else {
		return svc, err
	}
}
func ListSC() (*scv1.StorageClassList, error) {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	scs, err := clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	return scs, err
}
