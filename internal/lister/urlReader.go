package lister

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	ConnectionError = errors.New("could not fetch project from remote server")
)

type (
	urlReader struct {
		url    *url.URL
		client *http.Client
	}
)

func NewUrlReader(url *url.URL, client *http.Client) ReadStrategy {
	return &urlReader{
		url:    url,
		client: client,
	}
}

func (u *urlReader) ReadProjectStructures() (io.ReadCloser, error) {
	response, err := u.client.Get(u.url.String())

	if err != nil {
		return nil, ConnectionError
	}

	return response.Body, nil
}
