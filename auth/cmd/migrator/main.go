package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/IvanDrf/work-hunter/auth/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	migrations := ""

	flag.StringVar(&migrations, "mig", "", "path to migration files")
	cfg := config.LoadFromYAML()

	m, err := migrate.New(fmt.Sprintf("file://%s", migrations), cfg.Database.DSN())
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil {
		log.Fatal(err)
	}

}
