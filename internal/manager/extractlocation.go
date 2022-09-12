//+build !integration

package manager

import (
	"log"
	"os"
)

func GetExtractLocation() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("could not get working directory", err)
	}

	return wd
}
