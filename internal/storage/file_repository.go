package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"

	"github.com/AxMdv/go-url-shortener/internal/config"
)

// FileRepository is a file-based and in-memory-based repository.
type FileRepository struct {
	MapURL     map[string]string   // [shortened]long
	MapUUID    map[string][]string // [uuid][]shortened
	MapDeleted map[string]bool     // [shortened]deleted_flag
	filename   string
	URLSaver   *URLFileSaver
}

// NewFileRepository returns new FileRepository.
func NewFileRepository(config *config.Options) (*FileRepository, error) {
	repository := &FileRepository{
		MapURL:     make(map[string]string),
		MapUUID:    make(map[string][]string),
		MapDeleted: make(map[string]bool),
		filename:   config.FileStorage,
		URLSaver:   nil,
	}
	urlReader, err := NewURLReader(config.FileStorage)
	if err != nil {
		return nil, err
	}
	err = urlReader.ReadURL()
	if err != nil {
		return nil, err
	}
	for _, v := range urlReader.FormedURL {
		repository.MapURL[v.ShortenedURL] = v.LongURL
		repository.MapUUID[v.UUID] = append(repository.MapUUID[v.UUID], v.ShortenedURL)
	}
	urlReader.Close()
	urlSaver, err := NewURLFileSaver(config.FileStorage)
	if err != nil {
		return nil, err
	}
	repository.URLSaver = urlSaver

	return repository, nil
}

// AddURL writes url to FileRepository.
func (fr *FileRepository) AddURL(_ context.Context, formedURL *FormedURL) error {
	if fr.MapURL[formedURL.ShortenedURL] != "" {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	fr.MapURL[formedURL.ShortenedURL] = formedURL.LongURL

	fr.MapUUID[formedURL.UUID] = append(fr.MapUUID[formedURL.UUID], formedURL.ShortenedURL)

	err := fr.URLSaver.WriteURL(formedURL)
	return err
}

// AddURLBatch writes batch of urls to FileRepository.
func (fr *FileRepository) AddURLBatch(_ context.Context, formedURL []FormedURL) error {
	for _, v := range formedURL {
		fr.MapURL[v.ShortenedURL] = v.LongURL
		err := fr.URLSaver.WriteURL(&v)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetURL returns url from FileRepository.
func (fr *FileRepository) GetURL(_ context.Context, shortenedURL string) (string, error) {
	longURL := fr.MapURL[shortenedURL]
	if longURL == "" {
		return "", nil
	}
	return longURL, nil
}

// func (fr *FileRepository) GetURLByUserID(_ context.Context, uuid string) ([]FormedURL, error) {
// 	shortenedURL := fr.MapUUID[uuid]
// 	formedURL := make([]FormedURL, 0)
// 	for _, v := range shortenedURL {
// 		longURL, err := fr.GetURL(context.Background(), v)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var fu FormedURL
// 		fu.LongURL = longURL
// 		fu.ShortenedURL = v
// 		formedURL = append(formedURL, fu)
// 	}
// 	return formedURL, nil
// }

// GetURLByUserID returns urls shortened by user from FileRepository.
func (fr *FileRepository) GetURLByUserID(_ context.Context, uuid string) ([]FormedURL, error) {
	shortenedURL := fr.MapUUID[uuid]
	formedURL := make([]FormedURL, len(shortenedURL))
	for i, v := range shortenedURL {
		longURL, err := fr.GetURL(context.Background(), v)
		if err != nil {
			return nil, err
		}
		formedURL[i].LongURL = longURL
		formedURL[i].ShortenedURL = v
	}
	return formedURL, nil
}

// DeleteURLBatch deletes urls created by user.
func (fr *FileRepository) DeleteURLBatch(ctx context.Context, formedURL []FormedURL) error {
	for _, v := range formedURL {

		sliceShortened := fr.MapUUID[v.UUID]
		if sliceShortened == nil {
			continue
		}

		contains := contains(sliceShortened, v.ShortenedURL)
		if !contains {
			continue
		}
		fr.MapDeleted[v.ShortenedURL] = true
	}
	return nil
}

// GetFlagByShortURL returns if shortened url is deleted.
func (fr *FileRepository) GetFlagByShortURL(_ context.Context, shortenedURL string) (bool, error) {
	return fr.MapDeleted[shortenedURL], nil
}

// Close closes file of FileRepository.
func (fr *FileRepository) Close() error {
	err := fr.URLSaver.Close()
	return err
}

// URLFileSaver is object that provides writing urls data to file.
type URLFileSaver struct {
	file   *os.File
	writer *bufio.Writer
}

// NewURLFileSaver returns new URLFileSaver.
func NewURLFileSaver(filename string) (*URLFileSaver, error) {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &URLFileSaver{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// WriteURL writes FormedURL to file.
func (u *URLFileSaver) WriteURL(fu *FormedURL) error {
	data, err := json.Marshal(&fu)
	if err != nil {
		return err
	}
	if _, err := u.writer.Write(data); err != nil {
		return err
	}
	if err := u.writer.WriteByte('\n'); err != nil {
		return err
	}
	return u.writer.Flush()
}

// Close closes file containing url data.
func (u *URLFileSaver) Close() error {
	err := u.file.Close()
	return err
}

// URLFileReader is object that provides reading urls data from file.
type URLFileReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	FormedURL []FormedURL
}

// NewURLReader returns new URLFileReader.
func NewURLReader(filename string) (*URLFileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &URLFileReader{
		file:      file,
		scanner:   bufio.NewScanner(file),
		FormedURL: nil,
	}, nil
}

// ReadURL reads urls data from file.
func (u *URLFileReader) ReadURL() error {
	u.FormedURL = []FormedURL{}
	for u.scanner.Scan() {
		data := u.scanner.Bytes()
		tempFormed := FormedURL{}
		err := json.Unmarshal(data, &tempFormed)
		if err != nil {
			return err
		}
		u.FormedURL = append(u.FormedURL, tempFormed)
	}
	if err := u.scanner.Err(); err != nil {
		return u.scanner.Err()
	}
	return nil
}

// Close closes opened file.
func (u *URLFileReader) Close() error {
	return u.file.Close()
}
