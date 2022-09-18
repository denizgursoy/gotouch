<h1 align="center">gotouch</h1>
<p align="center">
<img alt="go report" src="https://goreportcard.com/badge/github.com/denizgursoy/gotouch"/>
<img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/denizgursoy/gotouch">
<a href="https://github.com/denizgursoy/gotouch/issues" target="_blank"><img alt="GitHub issues" src="https://img.shields.io/github/issues/denizgursoy/gotouch?color=b"></a>
</p>

gotouch easy way to create go projects. 

Install, run on CLI, make your selections and start development.

<p align="center">
<a href="https://asciinema.org/a/515980" target="_blank"><img alt="Asciiname" src="https://asciinema.org/a/515980.svg" /></a>
</p>

## Installation

```bash
go install -v github.com/denizgursoy/gotouch/cmd/gotouch@latest
```

### Features

- Choose from go project templates
- Easily install your dependencies and packages

# How to use
1. [Create your template project](#Create-your-template-project)
2. [Write your yaml file](#Write-your-yaml-file)


## Create your template project
Template is a zip file that has your directories and files inside.Template can be created with [package command](#package-subcommand).
Files inside a template can have [actions](https://pkg.go.dev/text/template#hdr-Actions) which will be [templated](https://pkg.go.dev/text/template)
with the [values](#values). As an example, If you have `Port` key in your [values](#values), `{{ .Port }}` will be replaced
with the corresponding value.

``` go
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", getRoot)
	log.Fatalln(http.ListenAndServe(":{{ .Port }}", nil))
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Server got the request\n")
}
```
You can also use other go template library's capabilities such as conditions, iterating array values etc. For more information
see [go template library](https://pkg.go.dev/text/template).

## Write your yaml file
Properties yaml is a **list** of what we call `Project Structure` as it can be seen below: 
```yaml
- name: Empty Project Layout #mandatory - used for displaying project in listing
  reference: https://go.dev/ #optional - is appended to name while prompting project name
  url: https://github.com/denizgursoy/go-touch-projects/raw/main/compressed/empty.zip #mandatory - url of template project 
  questions: #optional
    - direction: Do you want Dockerfile? #mandatory
      canSkip: true #if true, there must be at least one choice. 
      options:
        - choice: Yes
          files:
            - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
              pathFromRoot: Dockerfile
          values:
            isDocker: true
```
A Project Structure has name, reference, url and list of [questions](#Question) and described above.

After creating your template project, it should be hosted in http server and url must be a address of your template project. Gotouch downloads the template from the url 
and uncompress it. 


#### Question

### values

# Commands
### gotouch command

`gotouch --file path-to-yaml`

gotouch command uses [properties yaml](#Write-your-yaml-file) file for prompting user to enter name and select  project structure. If file flag value
is not provided, it is going to use [default properties yaml](https://github.com/denizgursoy/go-touch-projects/blob/main/properties.yaml). 
Firstly,The command asks for project name. Project name is written to go module and used for directory name.


### package subcommand
`gotouch package --source path-to-source --target path-to-target`

Package command compresses the source directory with the zip extension and moves the zip file to target directory.
`source` and `target` flags are optional. Default values for `source` and `target` are `./`, `../` respectively.


### validate subcommand
`gotouch validate --file path-to-yaml`

Validate checks if your yaml is valid or not. 

