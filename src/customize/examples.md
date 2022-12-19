# Examples

## Go Project

```yaml
- name: Empty Project Layout
  reference: https://go.dev/
§  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
```
## Projects in Other Languages

```yaml
- name: Empty Project Layout
  reference: https://maven.apache.org/
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
```
## Yes/No Question

```yaml
- name: Empty Project Layout
  reference: https://go.dev/
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  questions:
    - direction: Do you want Dockerfile? #mandatory
      canSkip: true #if true, there must be at least one choice.
      choices:
        - choice: Yes
          files:
            - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
              pathFromRoot: Dockerfile
```

## A Choice With All Fields
```yaml
- name: Empty Project Layout
  reference: https://go.dev/
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  questions:
    - direction: Do you want Dockerfile?
      canSkip: true
      choices:
        - choice: Yes
          dependencies:
            # dependency with latest
            - github.com/labstack/echo/v4
            # dependency with version
            - go.uber.org/zap@v1.23.0
          files:
            # file from internet
            - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
              pathFromRoot: Dockerfile
            # file from content
            - content: |
                lorem ipsum
                lorem ipsum
                lorem ipsum
                lorem ipsum
                lorem ipsum
              pathFromRoot: my-dir/input.txt # my-dir will be created if does not exist
          # values for templating
          # can be any type array, object, number, boolean and string
          values:
            #boolean
            isDocker: true
            # string
            myString: test
            # object
            myObject:
              x: value x
              y: value y
            # array
            myArray:
              - arrayValue1
              - arrayValue2
```
