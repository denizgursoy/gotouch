//go:build integration
// +build integration

package manager

import (
	"os"
)

func GetExtractLocation() string {
	return os.Getenv("TARGET_DIRECTORY")
}
