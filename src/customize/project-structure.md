# Project Structure

```yaml
- name: Empty Project Layout #mandatory 
  reference: https://go.dev/ #optional  
  url: https://github.com/denizgursoy/go-touch-projects.git # optional
  branch: empty # branch name of git repository to be cloned
  language: go # go, golang 
  delimiters: "<< >>" # optional to replace default templating values
  dependencies: # optional 
    - github.com/labstack/echo/v4
  files: # optional 
   - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
     pathFromRoot: Dockerfile
   - content: "foo"
     pathFromRoot: bar.txt
  values: # optional, cannot be changed by the user
    BaseURL: /v1
  customValues:  # optional, can be changed by the user
    Port: 8080
  questions: #optional
    - direction: Do you want Dockerfile? #mandatory
      canSkip: true #if true, there must be at least one choice. 
      choices:
        - choice: Yes
          files:
            - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
              pathFromRoot: Dockerfile
          values:
            isDocker: true
```

**`name`**: Name is displayed in project listing. It is a mandatory field.

**`reference`**: Reference is appended to name while prompting project's name

**`url`**: URL of template project. It can be a git repository `https://github.com/bff-project/bff.git`
or address of `.tar.gz` file like `https://github.com/denizgursoy/go-touch-projects/raw/main/compressed/empty.zip`. It
is a mandatory field.

**`branch`**: Gotouch allows you clone custom branch other than main on git. This allows you to host
all template projects in one git repository. This value is taken into consideration only if url is a git repository URL.

**`language`**: For Go project, it must be `go` or `golang`. For other languages, it can be omitted. Language is
appended to name while prompting project's name.

**`delimiters`**: Delimiter is an optional string field. It is used to replace go template library's default delimiter `{{`
and `}}`. New left and write delimiter should be seperated by space.

**`files`**: Same as files field of [choice](./choice.md). Allows you to create files when project structure is selected.

**`dependencies`**: Same as dependencies field of [choice](./choice.md). Allows you to add dependencies when project structure is selected.
It might be useful in some languages to get the latest version of a dependency.

**`values`**: Values used to templating files' content or directories' name in the template project. This filed can have
any values. It can be omitted as well. User cannot change these values.
See [value](./value) for more information.

**`customValues`**: Same as values but these values are prompted to user for change
See [value](./value) for more information.


**`questions`**: List of questions to customize project in case that this project structure is selected.
See [question](./question) for more information.

