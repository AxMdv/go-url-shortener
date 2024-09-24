package storage

import (
	"context"
	"database/sql"
	"errors"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = dbRepository.createDB(ctx)
	if err != nil {
		return nil, err
	}
	return &dbRepository, nil
}

func (dr *DBRepository) AddURL(ctx context.Context, formedURL *FormedURL) error {
	query := `
	INSERT INTO urls 
	VALUES ($1, $2, $3)
	ON CONFLICT ON CONSTRAINT urls_pk DO NOTHING;
	`
	result, err := dr.DB.ExecContext(ctx, query, formedURL.ShortenedURL, formedURL.LongURL, formedURL.UUID)
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

func (dr *DBRepository) AddURLBatch(ctx context.Context, formedURL []FormedURL) error {
	tx, err := dr.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls (shortened_url, long_url, uuid)"+
		" VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT urls_pk DO NOTHING")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, v := range formedURL {
		_, err := stmt.ExecContext(ctx, v.ShortenedURL, v.LongURL, v.UUID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (dr *DBRepository) GetURL(ctx context.Context, shortenedURL string) (string, error) {
	query := `
	SELECT long_url from urls WHERE shortened_url = $1;
	`
	rowLongURL := dr.DB.QueryRowContext(ctx, query, shortenedURL)
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

func (dr *DBRepository) createDB(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		shortened_url varchar NOT NULL,
		long_url varchar NOT NULL,
		uuid varchar NOT NULL,
		CONSTRAINT urls_pk PRIMARY KEY (shortened_url)
	);`
	_, err := dr.DB.ExecContext(ctx, query)
	return err
}

func (dr *DBRepository) PingDB(ctx context.Context, config config.Options) error {
	dsn := config.DataBaseDSN
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Panic(err)
		return err
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

func (dr *DBRepository) GetURLByUserID(ctx context.Context, uuid string) ([]FormedURL, error) {

	stmt, err := dr.DB.PrepareContext(ctx, "SELECT shortened_url, long_url FROM urls WHERE uuid = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultFormedURL := make([]FormedURL, 0)
	for rows.Next() {
		var fu FormedURL
		err := rows.Scan(&fu.ShortenedURL, &fu.LongURL)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, NewNoContentError(ErrNoContent, uuid)
			}
			return nil, err
		}
		resultFormedURL = append(resultFormedURL, fu)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return resultFormedURL, nil
}
