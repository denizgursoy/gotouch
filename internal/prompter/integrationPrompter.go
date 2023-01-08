//go:build integration

package prompter

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	urls  []string
	index = 0
)

func init() {
	getenv := os.Getenv("TARGET_FILE")
	body, err := os.ReadFile(getenv)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	urls = make([]string, 0)
	for _, line := range strings.Split(string(body), "\n") {
		urls = append(urls, line)
	}

}

func (s srv) AskForString(direction string, validator Validator) (string, error) {
	all, err := io.ReadAll(getStream())
	if err != nil {
		return "", err
	}
	return string(all), nil
}

func (s srv) AskForSelectionFromList(direction string, list []fmt.Stringer) (any, error) {
	all, err := io.ReadAll(getStream())
	if err != nil {
		return "", err
	}

	atoi, err := strconv.Atoi(string(all))
	return list[atoi], nil
}

func (s srv) AskForMultipleSelectionFromList(direction string, list []fmt.Stringer) ([]any, error) {
	all, err := io.ReadAll(getStream())
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(all))
	anies := make([]any, 0)
	for i, field := range fields {
		atoi, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}
		if atoi == 1 {
			anies = append(anies, list[i])
		}
	}
	return anies, nil
}

func (s srv) AskForYesOrNo(direction string) (bool, error) {
	all, err := io.ReadAll(getStream())
	if err != nil {
		return false, err
	}
	atoi, err := strconv.Atoi(string(all))

	return atoi == 1, nil
}

func (s srv) AskForMultilineString(direction, defaultValue, pattern string) (string, error) {
	return "", nil
}

func getStream() io.ReadCloser {
	nopCloser := io.NopCloser(strings.NewReader(urls[index]))
	index++
	return nopCloser
}
