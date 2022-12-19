# Template project 

Template project contains all the directories and files that you think you must have in every project.
Files inside a template can have [actions](https://pkg.go.dev/text/template#hdr-Actions) which will be [templated](https://pkg.go.dev/text/template)
with the [values](./value). For example, if you have an action like: 

 ```go
{{ .Port }} 
```
and `Port` key in your [values](./value), it  will be replaced with the corresponding value.

 ```go
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
Template project address is written to properties YAML. Template project can be in a git repository and the URL of `.git`
repository can be used in a properties YAML.

Template can be compressed as `.tar.gz` file with [package command](../commands#package). Compressed `.tar.gz` file can be hosted
in HTTP server and the URL to `.tar.gz` can be used in a properties YAML.

::: tip
The suggested way to store template projects is to host all template projects in one git repository. For every template
project, you can create orphan branches with `git checkout --orphan orphan_name` command and set `branch` in the project structure to `orphan_name`.
Common files can be stored in the `main` branch. See example projects: 

[Empty go project](https://github.com/denizgursoy/go-touch-projects/tree/empty)

[Standard go project](https://github.com/denizgursoy/go-touch-projects/tree/standard)

[Common files](https://github.com/denizgursoy/go-touch-projects)

To address files in the properties yaml, you can use [raw URL](https://www.howtogeek.com/wp-content/uploads/csit/2021/11/0ad2a42a.png?trim=1,1&bg-color=000&pad=1,1) of the files.
:::

 ### Templating with go templating library
You can also use other go template library's capabilities such as conditions, iterating array values, etc. For more
information see [go template library](https://pkg.go.dev/text/template).