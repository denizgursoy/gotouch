# How It Works

1. Asks the user to select a project structure in [properties YAML](../customize/properties-yaml) if there are more than
   on project structure in the YAML
2. If the selected project's language is go, it will check whether `go` command is installed on the OS
3. Asks for module name
4. Asks for a choice of every [question](../customize/question) under the selected project structure in order and saves
   the choices
5. Asks for changing the [values](../customize/values) of **only the selected project structure** (not selected choices)
   if any
6. Add default to [values](../customize/values)
7. Creates a new directory with module name's last part after last `/`
8. Uncompress/checkout the template project of the selected project structure into the created directory
9. If the selected project's language is go,updates the module's name in the go.mod with the value user entered, if
   there is no go.mod file, it creates the
   go.mod file
10. Creates files, and adds dependencies of all selected choices
11. Deletes `properties.yaml` on the root directory if exists.
12. If the selected project's language is go, executes `go mod tidy` and `go fmt ./...` on the root directory.
13. Merges values under the selected project structure with the values of all selected choices and default values
14. Walks through the newly created directory's content and templates every file with the
    merged [values](../customize/values)
15. Executes `init.sh`/`init.bat` on the root folder depending on the OS
16. Deletes `init.sh` and `init.bat` files on the root folder