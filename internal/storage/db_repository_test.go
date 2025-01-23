package storage

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/AxMdv/go-url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDBRepository(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}
	_, err := NewDBRepository(config)
	require.NoError(t, err)
}

func TestAddURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}

	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.DropTableURLS()
	require.NoError(t, err)

	formedURL := model.FormedURL{
		UUID:         "123",
		ShortenedURL: "228",
		LongURL:      "312",
	}
	db, err = NewDBRepository(config)
	require.NoError(t, err)
	err = db.AddURL(context.Background(), &formedURL)
	require.NoError(t, err)
	err = db.AddURL(context.Background(), &formedURL)
	var duplicateErr *AddURLError
	assert.True(t, errors.As(err, &duplicateErr))
}

func TestAddURLBatch(t *testing.T) {
	requestBatch := []model.FormedURL{
		{
			CorrelationID: "321",
			LongURL:       "https://vk.com",
			ShortenedURL:  "sssss",
		},
		{
			CorrelationID: "123",
			LongURL:       "https://yandex.ru",
			ShortenedURL:  "123",
		},
	}
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}

	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.DropTableURLS()
	require.NoError(t, err)
	db, err = NewDBRepository(config)
	require.NoError(t, err)
	err = db.AddURLBatch(context.Background(), requestBatch)
	require.NoError(t, err)
}

func TestGetURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}
	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.DropTableURLS()
	require.NoError(t, err)
	db, err = NewDBRepository(config)
	require.NoError(t, err)
	fu := model.FormedURL{
		CorrelationID: "321",
		LongURL:       "https://vk.com",
		ShortenedURL:  "sssss",
	}
	err = db.AddURL(context.Background(), &fu)
	require.NoError(t, err)

	long, err := db.GetURL(context.Background(), fu.ShortenedURL)
	require.NoError(t, err)
	assert.Equal(t, fu.LongURL, long)

	long, err = db.GetURL(context.Background(), "not existing key")
	assert.Equal(t, "", long)
	assert.NoError(t, err)
}

func TestPingDB(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}
	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.PingDB(context.Background())
	require.NoError(t, err)
}

func TestClose(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}
	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.Close()
	require.NoError(t, err)
}

func TestGetURLByUserID(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}
	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.DropTableURLS()
	require.NoError(t, err)
	db, err = NewDBRepository(config)
	require.NoError(t, err)

	fu := model.FormedURL{
		UUID:         "user-321",
		LongURL:      "https://vk.com",
		ShortenedURL: "sssss",
	}
	err = db.AddURL(context.Background(), &fu)
	require.NoError(t, err)

	formedResult, err := db.GetURLByUserID(context.Background(), fu.UUID)
	require.NoError(t, err)
	assert.Equal(t, fu.LongURL, formedResult[0].LongURL)
	assert.Equal(t, fu.ShortenedURL, formedResult[0].ShortenedURL)

	_, err = db.GetURLByUserID(context.Background(), "1232222")
	fmt.Println(err)
	var noContentError *NoContentError
	assert.True(t, errors.As(err, &noContentError))

}

func TestDeleteURLBatch_GetFlagByShortURL(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}
	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.DropTableURLS()
	require.NoError(t, err)
	db, err = NewDBRepository(config)
	require.NoError(t, err)

	fu := []model.FormedURL{{
		UUID:         "user-321",
		ShortenedURL: "sssss",
	}}
	err = db.AddURL(context.Background(), &fu[0])
	require.NoError(t, err)

	err = db.DeleteURLBatch(context.Background(), fu)
	require.NoError(t, err)

	t.Run("check deleted flag", func(t *testing.T) {
		deleted, err := db.GetFlagByShortURL(context.Background(), fu[0].ShortenedURL)
		require.NoError(t, err)
		assert.True(t, deleted)
		deleted, err = db.GetFlagByShortURL(context.Background(), "not deleted sth")
		require.NoError(t, err)
		assert.False(t, deleted)
	})

}

func TestGetNumberOfURLAndUser(t *testing.T) {
	config := &config.Options{
		RunAddr:            ":8080",
		ResponseResultAddr: "http://localhost:8080",
		FileStorage:        "",
		DataBaseDSN:        SetDSNForTests(),
	}
	db, err := NewDBRepository(config)
	require.NoError(t, err)
	err = db.DropTableURLS()
	require.NoError(t, err)
	db, err = NewDBRepository(config)
	require.NoError(t, err)
	_, err = db.GetNumberOfURLAndUser(context.Background())
	require.NoError(t, err)
}
