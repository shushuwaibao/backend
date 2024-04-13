package kubernetes

import (
	"context"
	"fmt"
	"gin-template/common"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func EditPVCSize(name string, size string) error {
	// 传入一个id和一个修改后的配置json，根据id修改配置
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
		pvc.Spec.Resources.Requests["storage"] = resource.MustParse(size)
		_, err := clientset.CoreV1().PersistentVolumeClaims("default").Update(context.TODO(), &pvc, metav1.UpdateOptions{})
		if err != nil {
			common.SysLog(err.Error())
		} else {
			common.SysLog(fmt.Sprintf("PVC %s size changed to %s", pvc.Name, size))
		}
	}
	return nil
}
