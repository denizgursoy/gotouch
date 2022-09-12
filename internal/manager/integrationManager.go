//+build integration

package manager

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetExtractLocation() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("could not fetch executable information", err)
	}
	fmt.Println("filepath.Dir(ex)", filepath.Dir(ex))
	return filepath.Dir(ex)
}
