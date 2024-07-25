package url

func (s *service) GetLongURL(shortened_URL string) (string, error) {
	longURL, err := s.urlRepository.GetURL(shortened_URL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}
