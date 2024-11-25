package k8s_client

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

type k8sClient interface {
	GetKubeConfig() (*Clientset, error)
	GetKubeEvents(*Clientset) (string, error)
}

type Clientset struct {
	client *kubernetes.Clientset
}

type EventsStruct struct {
	Name               string
	Namespace          string
	Count              int
	Type               string
	Event              string
	Reporter           string
	ObjectInvolvedKind string
	ObjectInvolvedName string
}

func GetKubeConfig() (*Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes clientset: %v", err)
	}

	return &Clientset{client: clientset}, err
}

func GetKubeEvetns(c *Clientset) string {
	eventsClient := c.client.CoreV1().Events("default")

	// List events with optional filtering
	eventsList, err := eventsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching events: %v", err)
	}

	var eventArray []EventsStruct
	// Iterate and print event details
	for _, event := range eventsList.Items {
		if event.Type == "Warning" {
			eventArray = append(eventArray, EventsStruct{
				Name:               event.Reason,
				Namespace:          event.ObjectMeta.Namespace,
				Count:              int(event.Count),
				Type:               event.Type,
				Event:              event.Reason, // Assuming Reason is the event message
				Reporter:           event.Source.Component,
				ObjectInvolvedKind: event.InvolvedObject.Kind,
				ObjectInvolvedName: event.InvolvedObject.Name,
			})
		}
	}
	fmt.Println(eventArray)

	return "got it"
}
