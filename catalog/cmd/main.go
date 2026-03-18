package main

import (
	"log"
	"os"
	"time"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/catalog"
	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	var repo catalog.Repository
	var err error

	cfg := elasticsearch.Config{
		CloudID: os.Getenv("ELASTIC_CLOUD_ID"),
		APIKey:  os.Getenv("ELASTIC_API_KEY"),
	}

	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		repo, err = catalog.NewElasticSearchRepository(cfg)
		if err == nil {
			log.Println("Connected to ElasticSearch successfully")
			break
		}

		log.Printf("Retry %d/%d: %s", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to ElasticSearch: %s", err)
	}

	service := catalog.NewService(repo)

	log.Println("Catalog gRPC service running on :50051 🚀")

	if err := catalog.ListenGRPC(service, ":50051"); err != nil {
		log.Fatalf("Failed to start gRPC server: %s", err)
	}
}
