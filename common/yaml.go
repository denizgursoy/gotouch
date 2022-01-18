package common

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func ReadYaml(myStruct interface{}, fileAddress string) error {
	file, err := ioutil.ReadFile(fileAddress)
	if err != nil {
		return nil
	}
	err = yaml.Unmarshal(file, myStruct)

	//fmt.Println(string(file))
	return err
}
