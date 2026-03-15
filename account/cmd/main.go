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
	service := account.NewAccountService(repo)
	account.ListenGRPC(service, cfg.GRPCPort)
	log.Printf("Account service is running on port %s", cfg.GRPCPort)
}
