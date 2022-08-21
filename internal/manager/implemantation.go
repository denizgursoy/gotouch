package manager

import (
	"fmt"
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

func newFileManager() Manager {
	return &fManager{}
}

func (f *fManager) CreateDirectoryIfNotExists(directoryName string) error {
	return os.Mkdir(directoryName, os.ModePerm)
}

func (f *fManager) GetStream() (ioReader io.ReadCloser) {
	if f.IsTest() {
		ioReader = io.NopCloser(strings.NewReader(urls[index]))
	} else {
		ioReader = os.Stdin
	}
	index++
	return
}

func (f *fManager) IsTest() bool {
	return Environment == "test"
}

func (f *fManager) GetExtractLocation() string {
	if f.IsTest() {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		return filepath.Dir(ex)
	} else {
		return f.GetWd()
	}
}

func (f *fManager) GetWd() string {
	getwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return getwd
}

func init() {
	manager = newFileManager()
	if manager.IsTest() {
		exPath := fmt.Sprintf("%s/input.txt", manager.GetExtractLocation())
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
