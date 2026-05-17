# Examples

### Go Project

```yaml
- name: Empty Project Layout
  reference: https://go.dev/
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
```
### Projects in Other Languages

```yaml
- name: Empty Project Layout
  reference: https://maven.apache.org/
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
```
### Yes/No Question

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

### A Choice With All Fields
```yaml
- name: Empty Project Layout
  reference: https://go.dev/
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  questions:
    - direction: Do you want Dockerfile?
      canSkip: true
      canSelectMultiple: false
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
            foo: test
            # object
            bar:
              x: value x
              y: value y
            # array
            fooBar:
              - arrayValue1
              - arrayValue2
```

### Docker + LocalPath

See [Distributing Templates with Docker](./local-path-docker-example) for a complete working example of the suggested
way to use Gotouch — building a self-contained Docker image that bundles your templates, configuration, and the tool
into a single versioned image your team can pull and run.
