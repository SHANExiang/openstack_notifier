package global

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"sincerecloud.com/openstack_notifier/consts"
)

func getNameSpace() string {
	content, err := ioutil.ReadFile(consts.NamespacePath)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func initK8sClient() *kubernetes.Clientset {
	// create Kubernetes config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	// create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func getK8sConfigMap() map[string]string {
	clientset := initK8sClient()
    namespace := getNameSpace()
    configMapName := "nonick-notifier-service"
    ctx := context.Background()
	configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	return configMap.Data
}

func GetServiceEnv() string {
	value := os.Getenv("ENV")

	// check env
	if value == "" {
		panic("Failed to get env value")
	}
	return value
}

func AssignCONF() {
	// get from k8s
	configMap := getK8sConfigMap()
	configMapKey := fmt.Sprintf("application-%s.yml", GetServiceEnv())
	content, ok := configMap[configMapKey]
	log.Println("content==", content)
	if !ok {
		panic("Failed to get configMap value")
	}
	if err := yaml.Unmarshal([]byte(content), &CONF); err != nil {
		panic(fmt.Sprintf("Failed to unmarshal configMap value %#v", err))
	}
	log.Println(fmt.Sprintf("CONF==%#v", CONF))
}

