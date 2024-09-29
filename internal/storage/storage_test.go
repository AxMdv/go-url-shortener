package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/AxMdv/go-url-shortener/internal/storage/mocks"
)

func TestStorageGetURL(t *testing.T) {
	type want struct {
		longURL string
	}
	tests := []struct {
		name         string
		expShort     string
		expLong      string
		shortenedURL string
		want         want
	}{
		{
			name:         "Positive test #1",
			expShort:     "aHR0cHM6Ly95YW5kZXgucnU",
			expLong:      "https://yandex.ru",
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				longURL: "https://yandex.ru",
			},
		},
		{
			name:         "Negative test #2",
			expShort:     "aHR0cHM6Ly95YW5kZXgucnU",
			expLong:      "",
			shortenedURL: "aHR0cHM6Ly95YW5kZXgucnU",
			want: want{
				longURL: "",
			},
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().GetURL(ctx, tt.shortenedURL).Return(tt.expLong, nil)
			longURL, err := m.GetURL(ctx, tt.shortenedURL)
			assert.Equal(t, tt.want.longURL, longURL)
			require.NoError(t, err)
		})
	}
}
