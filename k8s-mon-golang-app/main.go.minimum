package main

import (
    "fmt"
    "log"
    "os"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/tools/clientcmd"

    corev1 "k8s.io/api/core/v1"
)

func main() {
    kubeconfig := os.Getenv("KUBECONFIG")
    if kubeconfig == "" {
        log.Fatal("KUBECONFIG environment variable is not set")
    }

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Error building kubeconfig: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating Kubernetes client: %v", err)
    }

    // Create a ListWatch for Pods
    listWatch := cache.NewListWatchFromClient(
        clientset.CoreV1().RESTClient(),
        "pods",
        metav1.NamespaceAll,
        fields.Everything(),
    )

    informer := cache.NewSharedInformer(
        listWatch,
        &corev1.Pod{},
        0, // Resync period
    )

    informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            fmt.Println("Pod added")
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            fmt.Println("Pod updated")
        },
        DeleteFunc: func(obj interface{}) {
            fmt.Println("Pod deleted")
        },
    })

    stopCh := make(chan struct{})
    defer close(stopCh)
    informer.Run(stopCh)
}

