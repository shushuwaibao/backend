package kubernetes

import (
	"gin-template/common"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConf() *rest.Config {
	kubeconfig := os.Getenv("KUBE_CONFIG")
	tempconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		common.FatalLog(err)
	}
	return tempconfig
}
