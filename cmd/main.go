package main

import (
	"ClusterReport/internal/routes"
	"fmt"
	"net/http"
)

func main() {
	router := routes.NewRouter()

	//port := os.Getenv("PORT")
	port := 8080
	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
