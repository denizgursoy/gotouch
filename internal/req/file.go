package req

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/denizgursoy/gotouch/internal/logger"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/go-playground/validator/v10"
)

type (
	fileTask struct {
		File    model.File      `validate:"required"`
		Logger  logger.Logger   `validate:"required"`
		Manager manager.Manager `validate:"required"`
		Client  *http.Client    `validate:"required"`
	}
)

func (f *fileTask) Complete() error {
	if err := validator.New().Struct(f); err != nil {
		return err
	}
	url := f.File.Url
	var readCloser io.ReadCloser

	if len(strings.TrimSpace(url)) != 0 {
		resp, err := f.Client.Get(f.File.Url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		readCloser = resp.Body
	} else {
		reader := bytes.NewReader([]byte(f.File.Content))
		readCloser = io.NopCloser(reader)
	}

	path := f.File.PathFromRoot
	if len(strings.TrimSpace(path)) == 0 {
		path += filepath.Base(f.File.Url)
	}
	if err := f.Manager.CreateFile(readCloser, path); err != nil {
		return err
	}

	return nil
}
