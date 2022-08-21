package compressor

import (
	"fmt"
	"github.com/artdarek/go-unzip"
	"github.com/denizgursoy/gotouch/internal/manager"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type zipUncompressor struct {
	manager manager.Manager
}

func newZipUncompressor(manager manager.Manager) Compressor {
	return &zipUncompressor{
		manager: manager,
	}
}

func (z *zipUncompressor) UncompressFromUrl(url, projectName string) {
	client := http.Client{}
	response, httpErr := client.Get(url)
	if httpErr != nil {
		log.Fatalln("could connect to url", httpErr)
	}
	filePath := filepath.Join(os.TempDir(), filepath.Base(url))

	create, createFileErr := os.Create(filePath)
	if createFileErr != nil {
		log.Fatalln("could create file", httpErr)
	}

	if _, copyErr := io.Copy(create, response.Body); copyErr != nil {
		log.Fatalln("could not copy zip", copyErr)
	}
	target := fmt.Sprintf("%s/%s", z.manager.GetExtractLocation(), projectName)
	uz := unzip.New(filePath, target)

	if extractErr := uz.Extract(); extractErr != nil {
		log.Fatalln("could unzip the file", extractErr)
	}
}
