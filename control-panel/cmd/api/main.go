package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NirajDonga/neonpg/internal/handler"
	"github.com/NirajDonga/neonpg/internal/k8s"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	clientset, err := initKubernetesClient()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	log.Println("Successfully connected to Kubernetes cluster!")

	k8sProvisioner := k8s.NewProvisioner(clientset, "default")
	apiHandler := handler.NewAPIHandler(k8sProvisioner)

	http.HandleFunc("/create-database", apiHandler.HandleCreateDB)

	log.Println("Control Plane API running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return nil, homeErr
		}

		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
