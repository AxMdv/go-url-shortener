package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/AxMdv/go-url-shortener/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBRepository struct {
	DB *sql.DB
}

func NewDBRepository(config *config.Options) (*DBRepository, error) {
	db, err := sql.Open("pgx", config.DataBaseDSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	dbRepository := DBRepository{
		DB: db,
	}
	err = dbRepository.createDB()
	if err != nil {
		return nil, err
	}
	return &dbRepository, nil
}

func (dr *DBRepository) AddURL(formedURL *FormedURL) error {
	query := `
	INSERT INTO urls 
	VALUES ($1, $2, $3)
	ON CONFLICT ON CONSTRAINT urls_pk DO NOTHING;
	`
	result, err := dr.DB.ExecContext(context.Background(), query, formedURL.ShortenedURL, formedURL.LongURL, formedURL.UIID)
	if err != nil {
		return err
	}

	num, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if num == 0 {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	return nil
}

func (dr *DBRepository) AddURLBatch(formedURL []FormedURL) error {
	tx, err := dr.DB.Begin()
	if err != nil {
		return err
	}
	for _, v := range formedURL {
		_, err := tx.ExecContext(context.Background(), "INSERT INTO urls (shortened_url, long_url, uuid)"+
			" VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT urls_pk DO NOTHING", v.ShortenedURL, v.LongURL, v.UIID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return err
}

func (dr *DBRepository) GetURL(shortenedURL string) (string, error) {
	query := `
	SELECT long_url from urls WHERE shortened_url = $1;
	`
	rowLongURL := dr.DB.QueryRowContext(context.Background(), query, shortenedURL)
	var longURL string
	err := rowLongURL.Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (dr *DBRepository) Close() error {
	err := dr.DB.Close()
	return err
}

func (dr *DBRepository) createDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		shortened_url varchar NOT NULL,
		long_url varchar NOT NULL,
		uuid varchar NOT NULL,
		CONSTRAINT urls_pk PRIMARY KEY (shortened_url)
	);`
	_, err := dr.DB.ExecContext(context.Background(), query)
	return err
}

func (dr *DBRepository) PingDatabase(config config.Options) error {
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
