package util

import "net/url"

func EncodeURL(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	q := u.Query()
	u.RawQuery = q.Encode()
	return u.String(), nil
}
