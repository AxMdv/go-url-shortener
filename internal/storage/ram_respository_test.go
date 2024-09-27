package storage

import (
	"context"
	"testing"

	"math/rand"

	"github.com/stretchr/testify/assert"
)

func TestAddURL(t *testing.T) {
	type want struct {
		err error
	}
	tests := []struct {
		name          string
		formedURL     *FormedURL
		ramRepository RAMRepository
		want          want
	}{
		{
			name: "Positive test #1",
			formedURL: &FormedURL{
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

func BenchmarkAddURL(b *testing.B) {

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	src := rand.NewSource(123123123123)
	rand.New(src)
	b.Run("default", func(b *testing.B) {

		rr := RAMRepository{
			MapURL:     make(map[string]string),
			MapUUID:    make(map[string][]string),
			MapDeleted: make(map[string]bool),
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			randSeq1 := make([]rune, 20)
			randSeq2 := make([]rune, 20)
			for i := range randSeq1 {
				randSeq1[i] = letterRunes[rand.Intn(len(letterRunes))]
				randSeq2[i] = letterRunes[rand.Intn(len(letterRunes))]
			}
			randString1 := string(randSeq1)
			randString2 := string(randSeq2)

			rr.AddURL(context.Background(), &FormedURL{
				ShortenedURL: randString1,
				LongURL:      randString2,
			})

		}
	})
}
