package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	expserver "github.com/gmaschi/log-exp-eval/internal/servers/expressions"
	expstore "github.com/gmaschi/log-exp-eval/internal/services/datastore/postgresql/exp"
	"github.com/gmaschi/log-exp-eval/internal/services/eval"
	"github.com/gmaschi/log-exp-eval/pkg/tools/config/env"
	_ "github.com/lib/pq"
)

func main() {
	config, err := env.NewConfig()
	if err != nil {
		log.Printf("failed to load env variables: %v", err)
		os.Exit(1)
	}

	conn, err := sql.Open(
		config.DbDriver,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DBHost, config.DBPort, config.PostgresUser, config.PostgresPassword, config.PostgresDB),
	)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		os.Exit(1)
	}

	store := expstore.NewStore(conn)
	ev := eval.New()
	server, err := expserver.New(config, store, ev)
	if err != nil {
		log.Printf("failed to initialize server: %v", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
}
