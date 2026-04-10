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
	command := ""
	steps := 1

	flag.StringVar(&migrations, "mig", "", "path to migration files")
	flag.StringVar(&command, "cmd", "", "command for migration: up or down")
	flag.IntVar(&steps, "steps", 1, "amount of steps for migrations")
	cfg := config.LoadFromYAML()

	m, err := migrate.New(fmt.Sprintf("file://%s", migrations), cfg.Database.DSN())
	if err != nil {
		log.Fatal(err)
	}

	if command == "" {
		log.Fatalf("no command was given in cmd flags, need: up or down")
	}

	if command == "up" {
		if err = m.Steps(steps); err != migrate.ErrNoChange && err != nil {
			log.Fatal(err)
		}

		log.Println("UP migration applyed")
	}

	if command == "down" {
		if err = m.Steps(-steps); err != migrate.ErrNoChange && err != nil {
			log.Fatal(err)
		}

		log.Println("DOWN migration applyed")
	}
}
