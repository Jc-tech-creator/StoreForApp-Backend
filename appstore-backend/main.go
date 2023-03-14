package main

import (
	"fmt"
	"log"
	"net/http"

	"appstore/backend"
	"appstore/handler"
)

// main.go is the entrance of application

func main() {
    fmt.Println("started-service")
    backend.InitElasticsearchBackend()
	//This line initializes the Elasticsearch backend for the application.
	backend.InitGCSBackend()
	//This line initializes the Google Cloud Storage (GCS) backend for the application. 
    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}