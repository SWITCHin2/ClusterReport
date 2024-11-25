package routes

import (
	k8s_client "ClusterReport/internal/k8s-client"
	"fmt"
	"net/http"
	"sync"
)

type Route struct{}

var (
	clientSet  *k8s_client.Clientset
	clientOnce sync.Once
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/connectcluster", connectCluster)
	mux.HandleFunc("/event", getEvents)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "helath-check reached!")
}

func connectCluster(w http.ResponseWriter, r *http.Request) {
	clientOnce.Do(func() {
		var err error
		clientSet, err = k8s_client.GetKubeConfig()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create Kubernetes client: %s", err), http.StatusInternalServerError)
			return
		}
	})

	if clientSet != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Kubernetes client successfully initialized!")
	} else {
		http.Error(w, "Failed to initialize Kubernetes client", http.StatusInternalServerError)
	}
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	var _ = k8s_client.GetKubeEvetns(clientSet)
	// Write the YAML to the response
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
