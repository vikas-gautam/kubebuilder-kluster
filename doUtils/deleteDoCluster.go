package doutils

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	"github.com/digitalocean/godo"
// 	"k8s.io/client-go/kubernetes"
// )

// func Deletek8sCluster(c kubernetes.Interface, spec string, cname string) error {
// 	//tokenSecret has value in ns/secretname format
// 	secretNameSpace := strings.Split(spec.TokenSecret, "/")[0]
// 	k8sSecretName := strings.Split(spec.TokenSecret, "/")[1]

// 	//get the token value from k8sSecret dosecret
// 	tokenValue, err := getToken(secretNameSpace, k8sSecretName)
// 	if err != nil {
// 		fmt.Printf("Unable to get token from k8sSecret %s", err.Error())
// 	}

// 	//do client with tokenValue
// 	doClient := generateDoClient(tokenValue)
// 	fmt.Println(doClient)

// 	opt := &godo.ListOptions{
// 		Page:    1,
// 		PerPage: 200,
// 	}
// 	clusters, _, err := doClient.Kubernetes.List(context.Background(), opt)
// 	if err != nil {
// 		return err
// 	}

// 	for _, cluster := range clusters {
// 		if cluster.Name == cname {
// 			fmt.Println(cluster.ID)
// 			_, err = doClient.Kubernetes.Delete(context.Background(), cluster.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }
