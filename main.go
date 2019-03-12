package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
)

func main() {
	var ns, deployment, image, replicas string
	flag.StringVar(&ns, "namespace", "default", "namespace")
	flag.StringVar(&deployment, "deployment", "", "deployment")
	flag.StringVar(&image, "image", "", "image")
	flag.StringVar(&replicas, "replicas", "0", "replicas")
	flag.Parse()

	if *deployment == "" {
		fmt.Println("You must specify the deployment")
		os.Exit(0)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(ns)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var updateErr error

		result, err := deploymentsClient.Get(*deployment, metav1.Getoptions{})
		if err != nil {
			panic(err.Error())
		}
		if errors.IsNotFound(err) {
			fmt.Printf("Deployment not found\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting deployment%v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found deployment\n")
			name := result.GetName()
			fmt.Println("name ->", name)
			oldreplicas := *result.Spec.Replicas
			fmt.Println("old replicas ->", oldreplicas)
			fmt.Println("New replicas ->", *replicas)
			*result.Spec.Replicas = int32(*replicas)
			if *image != "" {
				*result.Spec.Template.Spec.Containers[0].Image = *image
			}
			_, updateErr := deploymentsClient.Update(result)
			fmt.Println(updateErr)
		}
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update deployment failed: %v", retryErr))
	}
	fmt.Println("Update deployment...")
}

