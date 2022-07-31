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

func (g gitHubLister) GetDefaultProjects() []*ProjectData {
	client := http.Client{}
	propertiesUrl := "https://raw.githubusercontent.com/denizgursoy/gotouch/main/projects/properties.yaml"
	response, err := client.Get(propertiesUrl)

	if err != nil {
		return nil
	}

	data := make([]*ProjectData, 0)

	allBytes, err := ioutil.ReadAll(response.Body)
	err = yaml.Unmarshal(allBytes, &data)

	return data
}
