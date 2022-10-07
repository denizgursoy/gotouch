package compressor

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
)

type (
	compressor struct {
		Manager  manager.Manager `validate:"required"`
		Store    store.Store     `validate:"required"`
		Strategy ZipStrategy     `validate:"required"`
	}
)

func newCompressor() Compressor {
	return &compressor{
		Manager:  manager.GetInstance(),
		Store:    store.GetInstance(),
		Strategy: newTarStrategy(),
	}
}

func (z *compressor) UncompressFromUrl(url string) error {
	if err := validator.New().Struct(z); err != nil {
		return err
	}

	client := http.Client{}
	response, httpErr := client.Get(url)
	if httpErr != nil {
		return httpErr
	}
	pattern := fmt.Sprintf("*%s", z.Strategy.GetExtension())
	temp, tempFileErr := os.CreateTemp("", pattern)
	if tempFileErr != nil {
		return tempFileErr
	}

	defer func() {
		os.Remove(temp.Name())
	}()

	if _, copyErr := io.Copy(temp, response.Body); copyErr != nil {
		return copyErr
	}

	projectName := z.Store.GetValue(store.ProjectName)
	target := fmt.Sprintf("%s/%s", z.Manager.GetExtractLocation(), projectName)

	return z.Strategy.UnCompressDirectory(temp.Name(), target)
}

func (z *compressor) CompressDirectory(source, target string) error {
	if !filepath.IsAbs(source) {
		absoluteSource, err := filepath.Abs(source)
		if err != nil {
			return err
		}
		source = absoluteSource
	}

	if !filepath.IsAbs(target) {
		absoluteTarget, err := filepath.Abs(target)
		if err != nil {
			return err
		}
		target = absoluteTarget
	}

	if !checkIsDirectory(source) {
		return errors.New("source is not a directory")
	} else if !checkIsDirectory(target) {
		return errors.New("target is not a directory")
	}

	if !strings.HasSuffix(source, string(os.PathSeparator)) {
		source = fmt.Sprintf("%s%s", source, string(os.PathSeparator))
	}

	filename := fmt.Sprintf("%s%s", filepath.Base(source), z.Strategy.GetExtension())
	target = filepath.Join(target, string(os.PathSeparator), filename)

	return z.Strategy.CompressDirectory(source, target)
}

func checkIsDirectory(path string) bool {
	open, err := os.Open(path)
	if err != nil {
		return false
	}
	stat, err := open.Stat()
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func shouldSkip(fileName string) bool {
	filesNotToZip := []string{"__MACOS", ".DS_Store", ".idea", ".vscode", ".git"}
	for _, file := range filesNotToZip {
		if strings.TrimSpace(fileName) == file {
			return true
		}
	}
	return false
}
