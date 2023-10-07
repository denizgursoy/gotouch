package compressor

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/auth"
	"github.com/denizgursoy/gotouch/internal/logger"

	"github.com/go-playground/validator/v10"

	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/store"
)

type (
	compressor struct {
		Client   *http.Client
		Manager  manager.Manager `validate:"required"`
		Store    store.Store     `validate:"required"`
		Strategy ZipStrategy     `validate:"required"`
		Logger   logger.Logger   `validate:"required"`
	}
)

func newCompressor() Compressor {
	return &compressor{
		Client:   auth.NewAuthenticatedHTTPClient(),
		Manager:  manager.GetInstance(),
		Store:    store.GetInstance(),
		Strategy: newTarStrategy(),
		Logger:   logger.NewLogger(),
	}
}

func (z *compressor) UncompressFromUrl(url string) error {
	if err := validator.New().Struct(z); err != nil {
		return err
	}

	z.Logger.LogInfo("Extracting files...")

	response, httpErr := z.Client.Get(url)
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

	if unCompressError := z.Strategy.UnCompressDirectory(temp.Name(), target); unCompressError != nil {
		return unCompressError
	}
	z.Logger.LogInfo("Zip is extracted successfully")

	return nil
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
