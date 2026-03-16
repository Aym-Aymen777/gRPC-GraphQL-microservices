package main

import (
	"log"

	"github.com/Aym-Aymen777/gRPC-GraphQL-microservices/account"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	MySQLURL string `env:"MYSQL_URL" envDefault:"root:password@tcp(localhost:3306)/accounts"`
	GRPCPort string `env:"GRPC_PORT" envDefault:":50051"`
}

func main() {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	repo, err := account.NewMySQLRepository(cfg.MySQLURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	log.Printf("Account service is connecting to MySQL ..... ⏳")
	service := account.NewAccountService(repo)
	log.Printf("Account service is connected to MySQL ✅")
	log.Printf("Account service is running on port %s 🎯", cfg.GRPCPort)
	account.ListenGRPC(service, cfg.GRPCPort)
}
