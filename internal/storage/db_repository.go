package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AxMdv/go-url-shortener/internal/config"
)

type DBRepository struct {
	db *pgxpool.Pool
}

func NewDBRepository(config *config.Options) (*DBRepository, error) {
	pool, err := pgxpool.New(context.Background(), config.DataBaseDSN)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	dbRepository := DBRepository{
		db: pool,
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
	result, err := dr.db.Exec(ctx, query, formedURL.ShortenedURL, formedURL.LongURL, formedURL.UUID)
	if err != nil {
		return err
	}

	num := result.RowsAffected()
	if num == 0 {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	return nil
}

func (dr *DBRepository) AddURLBatch(ctx context.Context, formedURL []FormedURL) error {
	tx, err := dr.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}
	query := "INSERT INTO urls (shortened_url, long_url, uuid)" +
		" VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT urls_pk DO NOTHING"

	for _, v := range formedURL {
		batch.Queue(query, v.ShortenedURL, v.LongURL, v.UUID)
	}
	dr.db.SendBatch(ctx, batch)
	return tx.Commit(ctx)
}

func (dr *DBRepository) GetURL(ctx context.Context, shortenedURL string) (string, error) {
	query := `
	SELECT long_url from urls WHERE shortened_url = $1;
	`
	rowLongURL := dr.db.QueryRow(ctx, query, shortenedURL)
	var longURL string
	err := rowLongURL.Scan(&longURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return "", nil
		}
		return "", err
	}
	return longURL, nil
}

func (dr *DBRepository) Close() error {
	dr.db.Close()
	return nil
}

func (dr *DBRepository) PingDB(ctx context.Context) error {
	err := dr.db.Ping(ctx)
	return err
}

func (dr *DBRepository) GetURLByUserID(ctx context.Context, uuid string) ([]FormedURL, error) {

	query := "SELECT shortened_url, long_url FROM urls WHERE uuid = $1"

	rows, err := dr.db.Query(ctx, query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultFormedURL := make([]FormedURL, 0)
	for rows.Next() {
		var fu FormedURL
		err := rows.Scan(&fu.ShortenedURL, &fu.LongURL)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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

func (dr *DBRepository) DeleteURLBatch(ctx context.Context, formedURL []FormedURL) error {

	query := `UPDATE urls SET deleted_flag = $1 WHERE uuid = $2 AND shortened_url = $3;`
	deletedFlag := true
	batch := &pgx.Batch{}
	for _, v := range formedURL {
		batch.Queue(query, deletedFlag, v.UUID, v.ShortenedURL)
	}

	br := dr.db.SendBatch(ctx, batch)
	defer br.Close()

	for range formedURL {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}
	return br.Close()

}

func (dr *DBRepository) GetFlagByShortURL(ctx context.Context, shortenedURL string) (bool, error) {
	query := `
	SELECT deleted_flag from urls WHERE shortened_url = $1;
	`
	row := dr.db.QueryRow(ctx, query, shortenedURL)

	var deleted bool
	err := row.Scan(&deleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return deleted, nil
}

func (dr *DBRepository) createDB(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS urls (
		shortened_url varchar NOT NULL,
		long_url varchar NOT NULL,
		uuid varchar NOT NULL,
		deleted_flag BOOLEAN NOT NULL DEFAULT FALSE,
		CONSTRAINT urls_pk PRIMARY KEY (shortened_url)
		);`
	_, err := dr.db.Exec(ctx, query)
	return err
}
