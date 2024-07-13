package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingDatabase() error {
	dsn := config.Options.DataBaseDSN
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Panic(err)
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Panic(err)
		return err
	}
	return nil
}
