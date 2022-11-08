package doutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	demov1alpha1 "github.com/vikas-gautam/kubebuilder-kluster/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getToken(secretNameSpace, k8sSecretName string) (string, error) {
	s, err := generateK8sClient().CoreV1().Secrets(secretNameSpace).Get(context.Background(), k8sSecretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// fmt.Println(s)
	return string(s.Data["token"]), nil
}

func Createk8sCluster(spec demov1alpha1.KlusterSpec) (string, error) {
	//tokenSecret has value in ns/secretname format
	secretNameSpace := strings.Split(spec.TokenSecret, "/")[0]
	k8sSecretName := strings.Split(spec.TokenSecret, "/")[1]

	//get the token value from k8sSecret dosecret
	tokenValue, err := getToken(secretNameSpace, k8sSecretName)
	if err != nil {
		fmt.Printf("Unable to get token from k8sSecret %s", err.Error())
	}

	//do client with tokenValue
	doClient := generateDoClient(tokenValue)
	fmt.Println(doClient)

	//call do clustercreation api
	request := &godo.KubernetesClusterCreateRequest{
		Name:        spec.Name,
		RegionSlug:  spec.Region,
		VersionSlug: spec.Version,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			&godo.KubernetesNodePoolCreateRequest{
				Name:  spec.NodePools[0].Name,
				Size:  spec.NodePools[0].Size,
				Count: spec.NodePools[0].Count,
			},
		},
	}

	clusterCreated, _, err := doClient.Kubernetes.Create(context.Background(), request)
	if err != nil {
		return "", err
	}

	return clusterCreated.ID, nil
}
