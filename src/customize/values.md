# Values

If you want some part of the source code not to be hardcoded, you can define custom values under
the `Project Structure`. The most common cases can be port numbers, service addresses, and some configuration values,
etc. Gotouch will ask user to change the values if he/she wants.

If following project structure is selected,

```yaml
- name: Empty Project Layout
  reference: https://go.dev/
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  values: # optional, cannot be changed by the user
    BaseURL: /v1
  customValues: # optional, can be changed by the user
    Port: 8080
```

Gotouch will ask whether user wants to change the **`custom values`** under the selected project structure.

![Edit Values](@images/edit-values.png)

The values are displayed in an inline text editor where you can modify them directly in the terminal.

![Vim editor](@images/vim-editor.png)

When you confirm, Gotouch will save the values and continue creating the project. Values under the selected
project structure will be merged with all selected choices' values.

::: warning
**`values`** of selected choices and the selected project structure cannot be changed. If you need some changeable
values
use  **`customValues`** of choices.
:::

### Default Values

Apart from these values, you can use following predefined values :

```yaml
ModuleName: Module name user typed (github.com/denigursoy/foo)
ProjectName: Project directory name (foo)
WorkingDirectory: Location where Gotouch command is executed (/tmp)
ProjectFullPath: Project's root directory (/tmp/foo)
Dependencies: Array of dependencies of all selected choices
```