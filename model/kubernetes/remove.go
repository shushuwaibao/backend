package kubernetes

import (
	"context"
	"fmt"
	"gin-template/common"
	"strings"

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
	err = clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		common.SysLog(err.Error())
	}
	return err
}

func RemoveService(name string, namespace string) error {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return err
	}
	err = clientset.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		common.SysLog(err.Error())
	}
	return err
}

func RemovePVC(name string) error {
	config := GetConf()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		common.SysLog(err.Error())
		return err
	}
	ssName := strings.Split(name, "-")[1]
	common.SysLog(fmt.Sprintf("PVC name: %s, SS name: %s", name, ssName))
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("default").List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=" + ssName,
	})

	// err = clientset.CoreV1().PersistentVolumeClaims("default").Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		common.SysLog(err.Error())
		return err
	}

	for _, pvc := range pvcs.Items {
		err = clientset.CoreV1().PersistentVolumeClaims("default").Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{})
		if err != nil {
			// common.SysLog(err.Error())
			// fmt.Println("Error deleting PVC:", pvc.Name, err)
			common.SysError(fmt.Sprintf("Error deleting PVC: %s, %v", pvc.Name, err))
		} else {
			common.SysLog(fmt.Sprintf("Successfully deleted PVC: %s", pvc.Name))
			// fmt.Printf("Successfully deleted PVC: %s\n", pvc.Name)
		}
	}
	return nil
}
