# Value

If you want some part of the source code not to be hardcoded, you can define custom values under the `Project Strcuture`
. The most common cases can be port numbers, service addresses, and some configuration values, etc. Gotouch will ask
user to change the values if he/she wants.

If following project structure is selected,

```yaml
- name: Empty Project Layout #mandatory 
  reference: https://go.dev/ #optional  
  url: https://github.com/denizgursoy/go-touch-projects/raw/main/compressed/empty.zip #mandatory 
  language: go # go, golang 
  values: # optional
    Port: 8080
```

Gotouch will ask whether user wants to change the values under the selected project structure.

![Edit Values](@images/edit-values.png)

Gotouch uses editor of [survey](https://github.com/go-survey/survey#editor). It launches your default editor for YAML.
If you want to change your editor, set **$VISUAL** or **$EDITOR** environment variables.

![Vim editor](@images/vim-editor.png)

When you exit your the editor, gotouch will save the values and continue creating the project. Values under the selected
project structure will be merged with all selected choices' values.

::: warning Be aware that only values under the selected project structure can be changed by users. Values of selected choices cannot be changed.
:::

## Default Values

Apart from these values, you can use following predefined values :

```
ModuleName: Module name user typed (github.com/denigursoy/my-app)
ProjectName: Project directory name (my-app)
WorkingDirectory: location where Gotouch command is executed
ProjectFullPath: Projects root directory
Dependencies: Array of dependencies of all selected choices