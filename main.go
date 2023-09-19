package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	outputFile, err := os.Create("deployments_without_limits.csv")
	if err != nil {
		panic(err.Error())
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// // Write the header
	// writer.Write([]string{"Namespace", "Deployment"})

	for _, namespace := range namespaces.Items {
		deployments, err := clientset.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Failed to list deployments in namespace %s: %v\n", namespace.Name, err)
			continue
		}

		for _, deployment := range deployments.Items {
			containers := deployment.Spec.Template.Spec.Containers
			for _, container := range containers {
				if container.Resources.Limits == nil {
					writer.Write([]string{deployment.Name, namespace.Name})
				}
			}
		}
	}

	fmt.Println("CSV file created!")
}
