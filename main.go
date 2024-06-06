package main

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// const (
// 	endpointURL = "http://your-endpoint-url" // Replace with your actual endpoint
// )

func main() {
	// Create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig file if running outside the cluster
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// Create a clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a list watcher for events
	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"events",
		metav1.Everything().String(),
		metav1.Everything(),
	)

	// Define the event handler
	_, controller := cache.NewInformer(
		watchlist,
		&corev1.Event{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				event := obj.(*corev1.Event)
				sendEvent(event)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				event := newObj.(*corev1.Event)
				sendEvent(event)
			},
		},
	)

	// Start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	// Wait forever
	select {}
}

func sendEvent(event *corev1.Event) {
	// Create a JSON payload from the event
	payload, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling event: %v\n", err)
		return
	}
	fmt.Printf("RJR received event payload from k8s ==> %v\n", payload)

	// Send the event to the endpoint
	// resp, err := http.Post(endpointURL, "application/json", bytes.NewBuffer(payload))
	// if err != nil {
	// 	fmt.Printf("Error sending event: %v\n", err)
	// 	return
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("Unexpected response code: %d\n", resp.StatusCode)
	// }
}
