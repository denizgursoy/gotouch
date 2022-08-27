package compressor

import (
	"fmt"
	"github.com/artdarek/go-unzip"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/model"
	"github.com/denizgursoy/gotouch/internal/store"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type zipCompressor struct {
	Manager manager.Manager
	Store   store.Store
}

func newZipCompressor(manager manager.Manager) Compressor {
	return &zipCompressor{
		Manager: manager,
		Store:   store.GetInstance(),
	}
}

func (z *zipCompressor) UncompressFromUrl(url string) error {
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

	projectName := z.Store.GetValue(store.ProjectName)
	target := fmt.Sprintf("%s/%s", z.Manager.GetExtractLocation(), projectName)
	uz := unzip.New(filePath, target)

	if extractErr := uz.Extract(); extractErr != nil {
		return extractErr
	}

	return nil
}

func isValid(compressor *zipCompressor) bool {
	if compressor.Manager == nil {
		return false
	}
	return true
}
