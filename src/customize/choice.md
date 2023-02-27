# Choice

If selected, a choice can create files, add dependencies and introduce new values.

Dependencies are a list. See [dependencies](./dependencies) section

A choice can create files with address of source file, or content. A file entry must have `pathFromRoot` value is the
location of the file inside project.

Creator of this yaml might want to customize project if a specific choice is selected, so values written under a choice
will be appended to general [values](./values). Values of choices cannot be changed by the user.

If you want user to change some values, you can write them under the **`customValues`**. Those values will be prompted
so that user can change.

A choice can be written like:

```yaml
- choice: Yes
  dependencies:
    - github.com/labstack/echo/v4
  files:
    - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
      pathFromRoot: Dockerfile
    - content: "foo"
      pathFromRoot: bar.txt
  values:
    httpLibrary: echo
  customValues:
    port: 8080
```