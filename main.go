package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    "time"

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
        log.Println("KUBECONFIG environment variable is not set. Using default kubeconfig.")
        kubeconfig = os.ExpandEnv("$HOME/.kube/config")
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
            pod, ok := obj.(*corev1.Pod)
            if !ok {
                log.Println("Error converting object to pod (AddFunc)")
                return
            }
            logPodDetails("Added", pod)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            pod, ok := newObj.(*corev1.Pod)
            if !ok {
                log.Println("Error converting object to pod (UpdateFunc)")
                return
            }
            logPodDetails("Updated", pod)
        },
        DeleteFunc: func(obj interface{}) {
            pod, ok := obj.(*corev1.Pod)
            if !ok {
                log.Println("Error converting object to pod (DeleteFunc)")
                return
            }
            logPodDetails("Deleted", pod)
        },
    })

    stopCh := make(chan struct{})
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
        }
    }()
    defer close(stopCh)

    log.Println("Starting informer...")
    informer.Run(stopCh)
}

// logPodDetails logs the detailed information of a pod
func logPodDetails(event string, pod *corev1.Pod) {
    nodeName := getSimpleNodeName(pod.Spec.NodeName)

    fmt.Printf(
        "[%s] Timestamp: %s, Namespace: %s, Pod Name: %s, Phase: %s, PodIP: %s, Node: %s, Restart Count: %d\n",
        event,
        pod.CreationTimestamp.Time.Format(time.RFC3339)[:19],
        pod.Namespace,
        pod.Name,
        pod.Status.Phase,
        pod.Status.PodIP,
        nodeName,
        calculateRestartCount(pod),
    )
    for _, containerStatus := range pod.Status.ContainerStatuses {
        state := "Unknown"
        if containerStatus.State.Waiting != nil {
            state = "Waiting - " + containerStatus.State.Waiting.Reason
        } else if containerStatus.State.Running != nil {
            state = "Running"
        } else if containerStatus.State.Terminated != nil {
            state = "Terminated - " + containerStatus.State.Terminated.Reason
        }
        fmt.Printf("  Container: %s, State: %s, Restarts: %d\n",
            containerStatus.Name, state, containerStatus.RestartCount)
    }
}

// calculateRestartCount calculates the total restarts for all containers in a pod
func calculateRestartCount(pod *corev1.Pod) int {
    totalRestarts := 0
    for _, containerStatus := range pod.Status.ContainerStatuses {
        totalRestarts += int(containerStatus.RestartCount)
    }
    return totalRestarts
}

// getSimpleNodeName extracts the part of the node name before the first dot
func getSimpleNodeName(nodeName string) string {
    if idx := strings.Index(nodeName, "."); idx != -1 {
        return nodeName[:idx]
    }
    return nodeName
}

