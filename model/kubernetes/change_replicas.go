package kubernetes

import (
	"context"
	"fmt"
	"gin-template/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ChangeReplicas(name string, namespace string, replicas int) error {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return err
	}
	common.SysLog(fmt.Sprintf("ns: %s, name: %s, replicas: %d", namespace, name, replicas))
	// pirnt all statefulsets:
	{
		ss, _ := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
		// if err != nil {
		// 	common.SysLog(err.Error())
		// 	return err
		for _, s := range ss.Items {
			common.SysLog(fmt.Sprintf("ns: %s, name: %s, replicas: %d", s.Namespace, s.Name, *s.Spec.Replicas))
		}
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
