package storage

import (
	"context"
	"os"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileRepository(t *testing.T) {
	config := &config.Options{
		FileStorage: "short-url-db.json",
	}
	fr, err := NewFileRepository(config)
	require.NoError(t, err)

	err = fr.Close()
	require.NoError(t, err)
	err = os.Remove("short-url-db.json")
	require.NoError(t, err)
}

func TestFileRepoAddURL(t *testing.T) {
	type want struct {
		err error
	}
	tests := []struct {
		name      string
		formedURL *FormedURL
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: &FormedURL{
				UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
				ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
				LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
			},
			want: want{
				err: nil,
			},
		},
	}
	config := &config.Options{
		FileStorage: "short-url-db.json",
	}
	fr, err := NewFileRepository(config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fr.AddURL(context.Background(), tt.formedURL)
			assert.Equal(t, tt.want.err, err)
			longURL := fr.MapURL[tt.formedURL.ShortenedURL]
			assert.Equal(t, tt.formedURL.LongURL, longURL)
		})
	}
	err = fr.Close()
	require.NoError(t, err)
	err = os.Remove("short-url-db.json")
	require.NoError(t, err)
}

func TestFileRepoAddURLBatch(t *testing.T) {
	type want struct {
		err error
	}
	tests := []struct {
		name      string
		formedURL []FormedURL
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: []FormedURL{
				{
					UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
					ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
					LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				},
			},
			want: want{
				err: nil,
			},
		},
	}
	config := &config.Options{
		FileStorage: "short-url-db.json",
	}
	fr, err := NewFileRepository(config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fr.AddURLBatch(context.Background(), tt.formedURL)
			assert.Equal(t, tt.want.err, err)
			for _, v := range tt.formedURL {
				assert.Contains(t, fr.MapURL, v.ShortenedURL)
			}
		})
	}
	err = fr.Close()
	require.NoError(t, err)
	err = os.Remove("short-url-db.json")
	require.NoError(t, err)
}

func TestFileRepoGetURL(t *testing.T) {
	type want struct {
		longURL string
		err     error
	}
	tests := []struct {
		name         string
		formedURL    *FormedURL
		shortenedURL string
		want         want
	}{
		{
			name: "Positive test #1",
			formedURL: &FormedURL{
				UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
				ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
				LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
			},
			shortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
			want: want{
				longURL: "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				err:     nil,
			},
		},
	}
	config := &config.Options{
		FileStorage: "short-url-db.json",
	}
	fr, err := NewFileRepository(config)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fr.AddURL(context.Background(), tt.formedURL)
			require.NoError(t, err)

			longURL, err := fr.GetURL(context.Background(), tt.shortenedURL)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.longURL, longURL)
		})
	}
	err = fr.Close()
	require.NoError(t, err)
	err = os.Remove("short-url-db.json")
	require.NoError(t, err)
}

func TestFileRepoGetURLByUserID(t *testing.T) {
	type want struct {
		formedURL []FormedURL
		err       error
	}
	tests := []struct {
		name      string
		uuid      string
		formedURL *FormedURL
		want      want
	}{
		{
			name: "Positive test #1",
			uuid: "01ef7cf6-286f-6e26-a782-00155dad7c8c",
			formedURL: &FormedURL{
				UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
				ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
				LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
			},
			want: want{
				formedURL: []FormedURL{
					{
						ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
						LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
					},
				},
				err: nil,
			},
		},
	}
	config := &config.Options{
		FileStorage: "short-url-db.json",
	}
	fr, err := NewFileRepository(config)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fr.AddURL(context.Background(), tt.formedURL)
			require.NoError(t, err)

			formedURL, err := fr.GetURLByUserID(context.Background(), tt.uuid)
			assert.Equal(t, tt.want.err, err)
			assert.ElementsMatch(t, tt.want.formedURL, formedURL)
		})
	}
	err = fr.Close()
	require.NoError(t, err)
	err = os.Remove("short-url-db.json")
	require.NoError(t, err)
}

func TestFileRepoDeleteURLBatch(t *testing.T) {
	type want struct {
		deletedFlag bool
		err         error
	}
	type maps struct {
		MapURL     map[string]string   //[shortened]long
		MapUUID    map[string][]string //[uuid][]shortened
		MapDeleted map[string]bool     //[shortened]deleted_flag
	}
	tests := []struct {
		name      string
		formedURL []FormedURL
		maps      maps
		want      want
	}{
		{
			name: "Positive test #1",
			formedURL: []FormedURL{
				{
					UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
					ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
					LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				},
			},
			maps: maps{
				MapURL: map[string]string{
					"aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1": "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				},
				MapUUID: map[string][]string{
					"01ef7cf6-286f-6e26-a782-00155dad7c8c": {"aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1"},
				},
				MapDeleted: make(map[string]bool),
			},
			want: want{
				deletedFlag: true,
				err:         nil,
			},
		},
	}
	config := &config.Options{
		FileStorage: "short-url-db.json",
	}
	fr, err := NewFileRepository(config)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr.MapURL = tt.maps.MapURL
			fr.MapUUID = tt.maps.MapUUID
			fr.MapDeleted = tt.maps.MapDeleted

			err := fr.DeleteURLBatch(context.Background(), tt.formedURL)
			assert.Equal(t, tt.want.err, err)
			for _, v := range tt.formedURL {
				deletedFlag := fr.MapDeleted[v.ShortenedURL]
				assert.Equal(t, tt.want.deletedFlag, deletedFlag)
			}
		})
	}
	err = fr.Close()
	require.NoError(t, err)
	err = os.Remove("short-url-db.json")
	require.NoError(t, err)
}

// func BenchmarkFileRepoAddURL(b *testing.B) {
// 	FileRepository := FileRepository
// }

func BenchmarkFileRepoGetURLByUserID(b *testing.B) {

	config := &config.Options{
		FileStorage: "short-url-db.json",
	}
	fr, err := NewFileRepository(config)
	require.NoError(b, err)

	userUUID := "01ef7cf6-286f-6e26-a782-00155dad7c8c"
	shortenedURL := []string{"aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL", "aHR0cDovL"}
	longURL := "http://sometext.com"
	fr.MapUUID[userUUID] = shortenedURL
	for _, v := range shortenedURL {
		fr.MapURL[v] = longURL
	}

	b.ResetTimer()
	b.Run("new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fr.GetURLByUserID(context.Background(), userUUID)
		}
	})
	fr.Close()
	require.NoError(b, err)
	err = os.Remove("short-url-db.json")
	require.NoError(b, err)
}
