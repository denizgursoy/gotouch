//go:build integration
// +build integration

package prompter

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	urls  []string
	index = 0
)

func init() {
	exPath := fmt.Sprintf("%s/input.txt", GetExtractLocation())
	file, err := os.ReadFile(exPath)
	if err != nil {
		log.Fatal(err)
	}
	urls = make([]string, 0)
	for _, line := range strings.Split(string(file), "\n") {
		urls = append(urls, line)
	}
}

func (s srv) AskForString(direction string, validator Validator) (string, error) {
	all, err := ioutil.ReadAll(getStream())
	if err != nil {
		return "", err
	}
	return string(all), nil

}

func (s srv) AskForSelectionFromList(direction string, list []fmt.Stringer) (interface{}, error) {
	all, err := ioutil.ReadAll(getStream())
	if err != nil {
		return "", err
	}

	atoi, err := strconv.Atoi(string(all))
	return list[atoi], nil
}

func (s srv) AskForYesOrNo(direction string) (bool, error) {
	all, err := ioutil.ReadAll(getStream())
	if err != nil {
		return false, err
	}
	atoi, err := strconv.Atoi(string(all))

	return atoi == 1, nil
}

func getStream() io.ReadCloser {
	nopCloser := io.NopCloser(strings.NewReader(urls[index]))
	index++
	return nopCloser
}

func GetExtractLocation() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("could not fetch executable information", err)
	}
	return filepath.Dir(ex)
}
