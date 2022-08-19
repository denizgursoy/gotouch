package prompts

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	urls        []string
	index       = 0
	Environment = "prod"
)

func init() {
	if IsTest() {
		exPath := fmt.Sprintf("%s/input.txt", GetExtractLocation())
		file, err := os.ReadFile(exPath)
		if err != nil {
			log.Println("deniz", err)
		}
		urls = make([]string, 0)
		for _, line := range strings.Split(string(file), "\n") {
			split := strings.Split(line, " ")
			ints := make([]byte, 0)

			for _, s := range split {
				atoi, _ := strconv.Atoi(s)
				ints = append(ints, byte(atoi))
			}
			urls = append(urls, string(ints))

		}
	}
}

type promptUi struct {
}

func (p promptUi) AskForSelectionFromList(direction string, listOptions []*ListOption) interface{} {
	options := make([]string, 0)
	for _, option := range listOptions {
		options = append(options, option.DisplayText)
	}

	prompt := promptui.Select{
		Label: direction,
		Items: options,
		Stdin: getStream(),
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.Println(err)
	}

	return listOptions[index].ReturnVal
}

func (p promptUi) AskForString(direction string, validator StringValidator) string {
	prompt := promptui.Prompt{
		Label:    direction,
		Validate: promptui.ValidateFunc(validator),
		Stdin:    getStream(),
	}
	run, err := prompt.Run()
	if err != nil {
		log.Println(err)
	}
	return run

}

func getStream() (ioReader io.ReadCloser) {
	if IsTest() {
		ioReader = io.NopCloser(strings.NewReader(urls[index]))
	} else {
		ioReader = os.Stdin
	}
	index++
	return
}

func IsTest() bool {
	return Environment == "test"
}

func GetExtractLocation() string {
	if IsTest() {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		return filepath.Dir(ex)
	} else {
		return GetWd()
	}
}

func GetWd() string {
	getwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return getwd
}
