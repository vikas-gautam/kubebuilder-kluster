package doutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	demov1alpha1 "github.com/vikas-gautam/kubebuilder-kluster/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getToken(secretNameSpace, k8sSecretName string) (string, error) {
	s, err := generateK8sClient().CoreV1().Secrets(secretNameSpace).Get(context.Background(), k8sSecretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// fmt.Println(s)
	return string(s.Data["token"]), nil
}

func getAllClusters(token string) ([]*godo.KubernetesCluster, error) {
	client := godo.NewFromToken(token)
	ctx := context.TODO()

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	clusters, _, err := client.Kubernetes.List(ctx, opt)

	if err != nil {
		return []*godo.KubernetesCluster{}, err
	}
	return clusters, nil

}

func Createk8sCluster(cl client.Client, kluster *demov1alpha1.Kluster) (string, error) {
	//tokenSecret has value in ns/secretname format
	secretNameSpace := strings.Split(kluster.Spec.TokenSecret, "/")[0]
	k8sSecretName := strings.Split(kluster.Spec.TokenSecret, "/")[1]

	//get the token value from k8sSecret dosecret
	tokenValue, err := getToken(secretNameSpace, k8sSecretName)
	if err != nil {
		fmt.Printf("Unable to get token from k8sSecret %s", err.Error())
	}

	//Get list of all existing clusters from DO
	allClusters, _ := getAllClusters(tokenValue)
	fmt.Printf("list of all existing clusters: %v", allClusters)

	//addfinalizer already updated the current kluster configuration with annotation (last applied) and finalizer field
	currentCluster, err := GetCurrentKluster(cl, kluster)
	if err != nil {
		log.Printf("error fetching existing cluster %s", err)
	}
	fmt.Println(currentCluster)

	//do client with tokenValue
	doClient := generateDoClient(tokenValue)
	fmt.Println(doClient)

	for _, cluster := range allClusters {

		if cluster.Name == currentCluster.Spec.Name {
			//if matches then call patchCluster func
			return cluster.Name, patchDOcluster(currentCluster, cluster, doClient)
		}
	}

	// call do clustercreation api

	request := &godo.KubernetesClusterCreateRequest{
		Name:        kluster.Spec.Name,
		RegionSlug:  kluster.Spec.Region,
		VersionSlug: kluster.Spec.Version,
		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			&godo.KubernetesNodePoolCreateRequest{
				Name:  kluster.Spec.NodePools[0].Name,
				Size:  kluster.Spec.NodePools[0].Size,
				Count: kluster.Spec.NodePools[0].Count,
			},
		},
	}

	clusterCreated, _, err := doClient.Kubernetes.Create(context.Background(), request)
	if err != nil {
		return "", err
	}

	return clusterCreated.ID, nil
}

func GetCurrentKluster(cl client.Client, kluster *demov1alpha1.Kluster) (*demov1alpha1.Kluster, error) {

	currentKluster := &demov1alpha1.Kluster{}

	err := cl.Get(context.Background(), types.NamespacedName{Name: kluster.Spec.Name, Namespace: kluster.Namespace}, currentKluster)
	if err != nil {
		return &demov1alpha1.Kluster{}, err
	}
	return currentKluster, nil
}

func patchDOcluster(currentKluster *demov1alpha1.Kluster, existingCluster *godo.KubernetesCluster, doClient *godo.Client) error {

	if currentKluster.Spec.NodePools[0].Count != existingCluster.NodePools[0].Count {
		return upgradeDONodePools(currentKluster.Spec.NodePools[0].Count, existingCluster.ID, existingCluster.NodePools[0], doClient)
	}
	return nil
}

func upgradeDONodePools(currentClusterNPCount int, ExistingClusterid string, np *godo.KubernetesNodePool, doClient *godo.Client) error {
	//code to update nodepoolsd
	token := os.Getenv("DIGITALOCEAN_TOKEN")

	client := godo.NewFromToken(token)
	ctx := context.TODO()

	upgradeRequest := &godo.KubernetesNodePoolUpdateRequest{
		Name:  np.Name,
		Count: &currentClusterNPCount,
	}
	_, _, err := client.Kubernetes.UpdateNodePool(ctx, ExistingClusterid, np.ID, upgradeRequest)
	if err != nil {
		return err
	}
	return nil
}
