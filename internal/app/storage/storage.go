package storage

import (
	"bufio"
	"encoding/json"
	"os"
)

type Repository struct {
	MapURL   map[string]string
	filename string
	URLSaver *URLSaver
}

func InitRepository(filename string) (*Repository, error) {
	// recoverURL := RecoverURL{}
	// err := recoverURL.Load("recover_url.json")
	// if err == nil
	// recoverURL := RecoverURL{Filename: filename}
	// recoverURL.SaveFilename(recoverURL.Filename)
	repository := &Repository{
		MapURL:   make(map[string]string),
		filename: filename,
		URLSaver: nil,
	}
	if filename == "" {
		return repository, nil
	}
	urlReader, err := NewURLReader(filename)
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
	urlSaver, err := NewURLSaver(filename)
	if err != nil {
		return nil, err
	}
	repository.URLSaver = urlSaver

	return repository, nil
}

type FormedURL struct {
	UIID         string `json:"uiid"`
	ShortenedURL string `json:"short_url"`
	LongURL      string `json:"original_url"`
}

func (r *Repository) AddURL(longURL string, shortenedURL string) {
	r.MapURL[shortenedURL] = longURL
}

func (r *Repository) FindShortenedURL(shortenedURL string) (string, bool) {
	longURL := r.MapURL[shortenedURL]
	if longURL == "" {
		return "", false
	}
	return longURL, true
}

/////////////////////////////////////////////////////////////////////////////////////////////

type URLSaver struct {
	file   *os.File
	writer *bufio.Writer
}

func NewURLSaver(filename string) (*URLSaver, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &URLSaver{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (u *URLSaver) WriteURL(fu *FormedURL) error {
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

type URLReader struct {
	file      *os.File
	scanner   *bufio.Scanner
	FormedURL []FormedURL
}

func NewURLReader(filename string) (*URLReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &URLReader{
		file:      file,
		scanner:   bufio.NewScanner(file),
		FormedURL: nil,
	}, nil
}

func (u *URLReader) ReadURL() error {
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

func (u *URLReader) Close() error {
	return u.file.Close()
}

// структура, запоминающая имя файла, в который сохранялись URL
// type RecoverURL struct {
//     Filename string `json:"recover_filename"`
// }

// func (r *RecoverURL) SaveFilename (fsname string) error {
// 	data, err := json.MarshalIndent(r, "", "	")
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(fsname, data, 0666)
// }

// func (r *RecoverURL) Load(fsname string) error {
//     data, err := os.ReadFile(fsname)
//     if err != nil {
//         return err
//     }
//     if err := json.Unmarshal(data, r); err != nil {
//         return err
//     }
//     return nil
// }
