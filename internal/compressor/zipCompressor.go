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
	temp, tempFileErr := os.CreateTemp("", "*.zip")
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
	uz := unzip.New(temp.Name(), target)

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
