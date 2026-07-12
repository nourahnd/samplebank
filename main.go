package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nourahnd/samplebank/api"
	db "github.com/nourahnd/samplebank/db/sqlc"
	"github.com/nourahnd/samplebank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db:%v", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
