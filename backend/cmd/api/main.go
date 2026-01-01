package main

import (
	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/database"
	"github.com/guatom999/self-boardcast/internal/server"
)

func main() {

	cfg := config.LoadConfig("../../.env")

	db := database.DatabaseConnect(cfg)

	server.NewServer(db, cfg).Start()

}
