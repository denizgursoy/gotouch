package compressor

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/artdarek/go-unzip"
	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/denizgursoy/gotouch/internal/store"
	"github.com/go-playground/validator/v10"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type zipCompressor struct {
	Manager manager.Manager `validate:"required"`
	Store   store.Store     `validate:"required"`
}

func newZipCompressor(manager manager.Manager) Compressor {
	return &zipCompressor{
		Manager: manager,
		Store:   store.GetInstance(),
	}
}

func (z *zipCompressor) UncompressFromUrl(url string) error {
	if err := validator.New().Struct(z); err != nil {
		return err
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

func (z *zipCompressor) CompressDirectory(path string) error {
	open, err := os.Open(path)
	if err != nil {
		return err
	}
	stat, err := open.Stat()
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("not a directory")
	} else {
		return zipDirectory(open.Name(), filepath.Base(open.Name()))
	}

	return nil
}

func zipDirectory(baseFolder, name string) error {

	// Get a Buffer to Write To
	outFile, err := os.Create(name + `.zip`)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}
