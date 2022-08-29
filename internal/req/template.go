package req

import (
	"github.com/denizgursoy/gotouch/internal/store"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type (
	templateTask struct {
		Store store.Store
	}
)

func (t *templateTask) Complete() error {
	path := t.Store.GetValue(store.ProjectFullPath)
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() != ".keep" {
				AddSimpleTemplate(path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func getArrayData() map[string]interface{} {
	data := make(map[string]interface{}, 4)
	data["Character"] = []string{"sky", "blue", "forest", "tavern", "cup", "cloud"}
	data["origin_year"] = 2019
	data["destination_year"] = 2052
	data["effect"] = "the world stability after the Apocalypse"
	return data
}

func getData() map[string]interface{} {
	data := make(map[string]interface{}, 4)
	data["character"] = "Jonas Kahnwald"
	data["origin_year"] = 2019
	data["destination_year"] = 2052
	data["effect"] = "the world stability after the Apocalypse"
	return data
}

func AddSimpleTemplate(path string) {
	files, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = files.Execute(f, getArrayData())
}

//func AddSimpleTemplate(path string) {
//	f, err := os.OpenFile(path, os.O_RDWR, 0755)
//	if err != nil {
//		log.Fatal(err)
//	}
//	all, err := ioutil.ReadAll(f)
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	err = f.Truncate(0)
//	_, err = f.Seek(0, 0)
//	if err != nil {
//		fmt.Println(err)
//	}
//	input := string(all)
//
//	tmpl := template.Must(template.New("email.tmpl").Parse(input))
//	err = tmpl.Execute(f, getData())
//	if err != nil {
//		panic(err)
//	}
//	if err := f.Close(); err != nil {
//		log.Fatal(err)
//	}
//}
