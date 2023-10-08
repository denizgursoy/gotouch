# Commands

### gotouch

gotouch command uses [properties yaml](./customize/properties-yaml) file for prompting user to enter name and select project
structure. If file flag value is not provided, it is going to
use [default properties yaml](https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/package.yaml).
Firstly,The command asks for project name. Project name is written to go module and used for directory name.

**`inline`**: Inline flag allows you to create projects in your current directory.

### package

`gotouch package --source path-to-source --target path-to-target`

Package command compresses the source directory with the .tar.gz extension and moves the zip file to target directory.
`source` and `target` flags are optional. Default values for `source` and `target` are `./`, `../` respectively.

Package command ignores following files/directories:

1. ***__MACOS***
2. ***.DS_Store***
3. ***.idea***
4. ***.vscode***
5. ***.git***

### validate

`gotouch validate --file path-to-yaml`

Validate checks if your yaml is valid or not. 

### config

`gotouch config`

Allows you to change following configurations:

**`url`**: Replaces the default URL. If changed, Gotouch will display project structures in the URL without `-f` flag.

Usage:

`gotouch config set url path-to-new-url` changes to default URL

`gotouch config unset url` removes the changed URL to default. 
