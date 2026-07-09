package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	dbDialect  = "postgres"
	driverName = "pgx"
)

func main() {
	var action string
	flag.StringVar(&action, "a", "", "action to migrations (up/down/reset)")
	flag.Parse()

	ctx := context.Background()

	db, err := sql.Open(driverName, os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to connect to pool: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	goose.SetDialect(dbDialect)

	dir := os.Getenv("MIGRATIONS_PATH")
	switch action {
	case "up":
		if err := goose.UpContext(ctx, db, dir); err != nil {
			fmt.Fprintf(os.Stdout, "failed to up migrations: %v", err)
			os.Exit(1)
		}
	case "down":
		if err := goose.DownContext(ctx, db, dir); err != nil {
			fmt.Fprintf(os.Stdout, "failed to down migrations: %v", err)
			os.Exit(1)
		}
	case "reset":
		if err := goose.ResetContext(ctx, db, dir); err != nil {
			fmt.Fprintf(os.Stdout, "failed to reset migrations: %v", err)
			os.Exit(1)
		}
	default:
		log.Fatalln("invalid migrate operation")
	}

	log.Println("migrations done")
}
