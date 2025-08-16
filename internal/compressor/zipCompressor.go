package compressor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/context/ctxhttp"

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

func (z *compressor) UncompressFromUrl(ctx context.Context, url string, directoryToCreateProject string) error {
	if err := validator.New().StructCtx(ctx, z); err != nil {
		return err
	}

	z.Logger.LogInfo("Extracting files...")

	response, err := ctxhttp.Get(ctx, z.Client, url)
	if err != nil {
		return err
	}

	if err = z.extract(response.Body, directoryToCreateProject); err != nil {
		return err
	}

	return nil
}

func (z *compressor) CheckIfFileExtensionIsSupported(source string) error {
	source, err := filepath.Abs(strings.TrimSpace(source))
	if err != nil {
		return err
	}

	if !strings.HasSuffix(source, z.Strategy.GetExtension()) {
		return fmt.Errorf("unexpexted file format, expexted format is (%s)", z.Strategy.GetExtension())
	}

	return nil
}

func (z *compressor) CopyDirectory(path string, directoryToCreateProject string) error {
	if err := os.CopyFS(directoryToCreateProject, os.DirFS(path)); err != nil {
		logger.NewLogger().LogErrorIfExists(err)
		return fmt.Errorf("could not copy the directory")
	}

	return nil
}

func (z *compressor) UncompressFromLocalPath(_ context.Context, source string, directoryToCreateProject string) error {
	localZipFile, err := os.Open(source)
	if err = z.extractWithoutCopyingToTemp(localZipFile, directoryToCreateProject); err != nil {
		return err
	}

	return nil
}

func (z *compressor) CompressDirectory(source, target string) error {
	source, err := filepath.Abs(strings.TrimSpace(source))
	if err != nil {
		return err
	}

	target, err = filepath.Abs(strings.TrimSpace(target))
	if err != nil {
		return err
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

func (z *compressor) extract(fileReader io.Reader, directoryToCreateProject string) error {
	pattern := fmt.Sprintf("*%s", z.Strategy.GetExtension())
	temp, err := os.CreateTemp("", pattern)
	if err != nil {
		return err
	}

	defer func() {
		err := os.Remove(temp.Name())
		z.Logger.LogErrorIfExists(err)
	}()

	if _, err = io.Copy(temp, fileReader); err != nil {
		z.Logger.LogErrorIfExists(err)

		return err
	}

	if err = z.extractWithoutCopyingToTemp(temp, directoryToCreateProject); err != nil {
		z.Logger.LogErrorIfExists(err)

		return err
	}

	return nil
}

func (z *compressor) extractWithoutCopyingToTemp(compressedFile *os.File, directoryToCreateProject string) error {
	if err := z.Strategy.UnCompressDirectory(compressedFile.Name(), directoryToCreateProject); err != nil {
		return err
	}
	z.Logger.LogInfo("Zip is extracted successfully")

	return nil
}
