package lister

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type gitHubLister struct {
}

func newGithubLister() Lister {
	return gitHubLister{}
}

func (g gitHubLister) GetDefaultProjects() []*ProjectStructureData {
	client := http.Client{}
	propertiesUrl := "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/properties.yaml"
	response, err := client.Get(propertiesUrl)

	if err != nil {
		return nil
	}

	data := make([]*ProjectStructureData, 0)

	allBytes, err := ioutil.ReadAll(response.Body)
	err = yaml.Unmarshal(allBytes, &data)

	return data
}
