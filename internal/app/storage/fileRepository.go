package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/AxMdv/go-url-shortener/internal/app/config"
)

type FileRepository struct {
	MapURL   map[string]string
	filename string
	URLSaver *URLFileSaver
}

func (fr *FileRepository) AddURL(formedURL *FormedURL) error {
	if fr.MapURL[formedURL.ShortenedURL] != "" {
		return NewDuplicateError(ErrDuplicate, formedURL.ShortenedURL)
	}
	fr.MapURL[formedURL.ShortenedURL] = formedURL.LongURL

	err := fr.URLSaver.WriteURL(formedURL)
	return err
}

func (fr *FileRepository) AddURLBatch(formedURL []FormedURL) error {
	for _, v := range formedURL {
		fr.MapURL[v.ShortenedURL] = v.LongURL
		err := fr.URLSaver.WriteURL(&v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fr *FileRepository) GetURL(shortenedURL string) (string, bool) {
	longURL := fr.MapURL[shortenedURL]
	if longURL == "" {
		return "", false
	}
	return longURL, true
}

func (fr *FileRepository) Close() error {
	err := fr.URLSaver.Close()
	return err
}

func NewFileRepository(config *config.Options) (*FileRepository, error) {
	repository := &FileRepository{
		MapURL:   make(map[string]string),
		filename: config.FileStorage,
		URLSaver: nil,
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
	}
	urlReader.Close()
	urlSaver, err := NewURLFileSaver(config.FileStorage)
	if err != nil {
		return nil, err
	}
	repository.URLSaver = urlSaver

	return repository, nil
}

type URLFileSaver struct {
	file   *os.File
	writer *bufio.Writer
}

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

func (u *URLFileSaver) Close() error {
	err := u.file.Close()
	return err
}

type URLFileReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	FormedURL []FormedURL
}

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

func (u *URLFileReader) Close() error {
	return u.file.Close()
}
