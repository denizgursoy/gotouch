package compressor

import (
	"fmt"
	"github.com/artdarek/go-unzip"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type zipCompressor struct {
	manager manager.Manager
}

func newZipCompressor(manager manager.Manager) Compressor {
	return &zipCompressor{
		manager: manager,
	}
}

func (z *zipCompressor) UncompressFromUrl(url, projectName string) error {
	if !isValid(z) {
		return model.ErrMissingField
	}

	client := http.Client{}
	response, httpErr := client.Get(url)
	if httpErr != nil {
		return httpErr
	}
	filePath := filepath.Join(os.TempDir(), filepath.Base(url))

	create, createFileErr := os.Create(filePath)
	if createFileErr != nil {
		return createFileErr
	}

	if _, copyErr := io.Copy(create, response.Body); copyErr != nil {
		return copyErr
	}
	target := fmt.Sprintf("%s/%s", z.manager.GetExtractLocation(), projectName)
	uz := unzip.New(filePath, target)

	if extractErr := uz.Extract(); extractErr != nil {
		return extractErr
	}

	return nil
}

func isValid(compressor *zipCompressor) bool {
	if compressor.manager == nil {
		return false
	}
	return true
}
