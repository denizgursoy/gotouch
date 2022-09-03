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
	"strings"
)

const (
	compressExtension = ".zip"
)

type (
	zipCompressor struct {
		Manager manager.Manager `validate:"required"`
		Store   store.Store     `validate:"required"`
	}
)

func newZipCompressor() Compressor {
	return &zipCompressor{
		Manager: manager.GetInstance(),
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
	pattern := fmt.Sprintf("*.%s", compressExtension)
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
	uz := unzip.New(temp.Name(), target)

	if extractErr := uz.Extract(); extractErr != nil {
		return extractErr
	}

	return nil
}

func (z *zipCompressor) CompressDirectory(source, target string) error {
	if len(strings.TrimSpace(target)) == 0 {
		getwd, _ := os.Getwd()
		target = getwd
	}

	if !checkIsDirectory(source) {
		return errors.New("source is not a directory")
	} else if !checkIsDirectory(target) {
		return errors.New("target is not a directory")
	}

	if !strings.HasSuffix(source, string(os.PathSeparator)) {
		source = fmt.Sprintf("%s%s", source, string(os.PathSeparator))
	}

	filename := fmt.Sprintf("%s%s", filepath.Base(source), compressExtension)
	join := filepath.Join(target, string(os.PathSeparator), filename)

	return zipDirectory(source, join)
}

func zipDirectory(sourceFolder, targetFilePath string) error {
	outFile, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, sourceFolder, "")

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
			//	fmt.Println("Recursing and Adding SubDir: " + file.Name())
			//		fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
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
