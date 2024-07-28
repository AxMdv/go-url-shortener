package storage

import (
	"context"
	"database/sql"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
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
	_, err := dr.DB.ExecContext(context.Background(), query, formedURL.ShortenedURL, formedURL.LongURL, formedURL.UIID)
	if err != nil {
		return err
	}
	return nil
}

func (dr *DBRepository) GetURL(shortenedURL string) (string, bool) {
	query := `
	SELECT long_url from urls WHERE shortened_url = $1;
	`
	rowLongURL := dr.DB.QueryRowContext(context.Background(), query, shortenedURL)
	var longURL string
	err := rowLongURL.Scan(&longURL)
	if err != nil {
		return "", false
	}
	return longURL, true
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
