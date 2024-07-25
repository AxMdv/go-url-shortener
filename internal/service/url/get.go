package url

func (s *service) GetLongURL(shortenedURL string) (string, error) {
	longURL, err := s.urlRepository.GetURL(shortenedURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}
