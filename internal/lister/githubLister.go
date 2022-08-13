package lister

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ConnectionError = errors.New("could not fetch project from remote server")
	PropertiesUrl   = "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/properties.yaml"
)

type httpLister struct {
	client *http.Client
	URL    *string
}

func newHttpLister(client *http.Client, url *string) Lister {
	return httpLister{
		client: client,
		URL:    url,
	}
}

func (h httpLister) GetDefaultProjects() ([]*ProjectStructureData, error) {

	response, err := h.client.Get(*h.URL)

	if err != nil {
		log.Println(err.Error())
		return nil, ConnectionError
	}

	data := make([]*ProjectStructureData, 0)

	allBytes, err := ioutil.ReadAll(response.Body)
	err = yaml.Unmarshal(allBytes, &data)

	//TODO: data == nil durumunu handle et

	if data == nil {
		return data, errors.New("data cannot be empty")
	}

	return data, nil
}
