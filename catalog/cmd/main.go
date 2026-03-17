package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	cfg := elasticsearch.Config{
		CloudID: "CLOUD_ID",
		APIKey:  "API_KEY",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
}
