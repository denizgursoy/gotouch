package lister

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type gitHubLister struct {
}

func newGithubLister() Lister {
	return gitHubLister{}
}

func (g gitHubLister) GetDefaultProjects() ([]*ProjectStructureData, error) {
	client := http.Client{}
	propertiesUrl := "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/properties.yaml"
	response, err := client.Get(propertiesUrl)

	if err != nil {
		return nil, err
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
