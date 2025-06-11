package main

import (
	"flag"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"path/filepath"

	pdb "JWTService/pkg/postgres"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	up := flag.Bool("up", false, "Run up migrations")
	down := flag.Bool("down", false, "Run down migrations")
	flag.Parse()

	if !*up && !*down {
		log.Fatal("Please specify either -up or -down flag")
	}

	db, err := pdb.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("Failed to get absolute path to migrations: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	if *up {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run up migrations: %v", err)
		}
		log.Println("Successfully ran up migrations")
	}

	if *down {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run down migrations: %v", err)
		}
		log.Println("Successfully ran down migrations")
	}
}
