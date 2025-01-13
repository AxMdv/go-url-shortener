package storage

import (
	"context"
	"testing"

	"github.com/AxMdv/go-url-shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRAMRepoAddURL(t *testing.T) {
	type want struct {
		err error
	}
	tests := []struct {
		name          string
		formedURL     *model.FormedURL
		ramRepository RAMRepository
		want          want
	}{
		{
			name: "Positive test #1",
			formedURL: &model.FormedURL{
				UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
				ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
				LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
			},
			ramRepository: RAMRepository{
				MapURL:     make(map[string]string),
				MapUUID:    make(map[string][]string),
				MapDeleted: make(map[string]bool),
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.ramRepository.AddURL(context.Background(), tt.formedURL)
			assert.Equal(t, tt.want.err, err)
			longURL := tt.ramRepository.MapURL[tt.formedURL.ShortenedURL]
			assert.Equal(t, tt.formedURL.LongURL, longURL)
		})
	}
}

func BenchmarkRAMRepoAddURL(b *testing.B) {
	rr := RAMRepository{
		MapURL:     make(map[string]string),
		MapUUID:    make(map[string][]string),
		MapDeleted: make(map[string]bool),
	}
	formedURL := &model.FormedURL{
		UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
		ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
		LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
	}
	b.ResetTimer()
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rr.AddURL(context.Background(), formedURL)
		}
	})
}

func TestRAMRepoGetURL(t *testing.T) {
	type want struct {
		longURL string
		err     error
	}
	tests := []struct {
		name          string
		shortenedURL  string
		ramRepository RAMRepository
		want          want
	}{
		{
			name:         "Positive test #1",
			shortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
			ramRepository: RAMRepository{
				MapURL: map[string]string{
					"aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1": "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				},
				MapUUID:    make(map[string][]string),
				MapDeleted: make(map[string]bool),
			},
			want: want{
				longURL: "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			longURL, err := tt.ramRepository.GetURL(context.Background(), tt.shortenedURL)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.longURL, longURL)
		})
	}
}

func BenchmarkRAMRepoGetURL(b *testing.B) {
	rr := RAMRepository{
		MapURL:     make(map[string]string),
		MapUUID:    make(map[string][]string),
		MapDeleted: make(map[string]bool),
	}
	formedURL := &model.FormedURL{
		UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
		ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
		LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
	}
	b.ResetTimer()
	b.Run("default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rr.GetURL(context.Background(), formedURL.ShortenedURL)
		}
	})
}
func TestRAMRepoGetURLByUserID(t *testing.T) {
	type want struct {
		formedURL []model.FormedURL
		err       error
	}
	tests := []struct {
		name          string
		uuid          string
		ramRepository RAMRepository
		want          want
	}{
		{
			name: "Positive test #1",
			uuid: "01ef7cf6-286f-6e26-a782-00155dad7c8c",
			ramRepository: RAMRepository{
				MapURL: map[string]string{
					"aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1": "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				},
				MapUUID: map[string][]string{
					"01ef7cf6-286f-6e26-a782-00155dad7c8c": {"aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1"},
				},
				MapDeleted: make(map[string]bool),
			},
			want: want{
				formedURL: []model.FormedURL{
					{
						ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
						LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
					},
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			formedURL, err := tt.ramRepository.GetURLByUserID(context.Background(), tt.uuid)
			assert.Equal(t, tt.want.err, err)
			assert.ElementsMatch(t, tt.want.formedURL, formedURL)
		})
	}
}

func TestRAMRepoDeleteURLBatch(t *testing.T) {
	type want struct {
		deletedFlag bool
		err         error
	}
	tests := []struct {
		name          string
		formedURL     []model.FormedURL
		ramRepository RAMRepository
		want          want
	}{
		{
			name: "Positive test #1",
			formedURL: []model.FormedURL{
				{
					UUID:         "01ef7cf6-286f-6e26-a782-00155dad7c8c",
					ShortenedURL: "aHR0cDovL2ZwMXZpZTh0dXphMXB0LnJ1L3Ztd3dxMndhL2xuYjR1",
					LongURL:      "http://fp1vie8tuza1pt.ru/vmwwq2wa/lnb4u",
				},
			},
			ramRepository: RAMRepository{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.ramRepository.DeleteURLBatch(context.Background(), tt.formedURL)
			assert.Equal(t, tt.want.err, err)
			for _, v := range tt.formedURL {
				deletedFlag := tt.ramRepository.MapDeleted[v.ShortenedURL]
				assert.Equal(t, tt.want.deletedFlag, deletedFlag)
			}
		})
	}
}
