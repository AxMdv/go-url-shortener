package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func PingDatabase(config config.Options) error {
	dsn := config.DataBaseDSN
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
