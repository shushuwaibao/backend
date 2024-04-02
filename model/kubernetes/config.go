package kubernetes

import (
	"gin-template/common"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var conf *rest.Config

func GetConf() *rest.Config {
	if conf != nil {
		return conf
	}
	kubeconfig := os.Getenv("KUBE_CONFIG")
	tempconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// fmt.Print(conf.id)
	if err != nil {
		common.FatalLog(err)
	}
	conf = tempconfig
	return conf
}
