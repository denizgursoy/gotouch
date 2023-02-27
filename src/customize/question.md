# Question

Question allows your users to customize their projects. A question must have a direction and choices. Gotouch prompts
user to make a [choice](./choice) for every question in the selected Project Structure.

**`direction`**: Direction is displayed on the screen for the question. It is a mandatory field.

**`canSkip`**: This field allows skipping question. If there is one choice, question will be displayed as [yes/no question](#yes-no-question).
If there are more than one choice, it will be displayed as [None of above question](#none-of-above-question). It is optional field.

**`canSelectMultiple`**: This field allows users to make more than one choice. Instead of asking many yes/no questions, choices can be combined
in one question. Question will be displayed as [multiple select question](#multiple-select-question). It is optional field.

**`choices`**: Choices filed is a list of choices belonging to the question. It is mandatory field. See [choice](./choice) for more information.

### Yes/No question
If a question has only one choice and `canSkip` is true, it is evaluated as Yes/No question

<code-group>
<code-block title="Terminal">

![Yes/No Question](@images/yes-no-question.png)

</code-block>

<code-block title="YAML">

```yaml
- name: Api Gateway
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  questions:
    - direction: Do you want Dockerfile?
      canSkip: true #if true, there must be at least one choice. 
      choices: #mandatory
        - choice: Yes # mandatory
          files:
            - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
              pathFromRoot: Dockerfile # mandatory
```

</code-block>
</code-group>

### Multiple choice question

If a question has more than one choice it is evaluated as Multiple choice question

<code-group>
<code-block title="Terminal">

![Multiple Choice Question](@images/multiple-choice.png)

</code-block>

<code-block title="YAML">

```yaml
- name: Api Gateway
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  questions:
    - direction: Which HTTP framework do you want to use?
      choices:
        - choice: Echo
          dependencies:
            - github.com/labstack/echo/v4
        - choice: Gorilla Mux
          dependencies:
            - github.com/gorilla/mux
        - choice: Gin
          dependencies:
            - github.com/gin-gonic/gin
```

</code-block>
</code-group>


### None of above question

If a question has more than one choice and `canSkip` is true, `None of above` option will be also added 
as a choice. 

<code-group>
<code-block title="Terminal">

![None of Above Question](@images/none-of-above.png)

</code-block>

<code-block title="YAML">

```yaml
- name: Api Gateway
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  questions:
    - direction: Which HTTP framework do you want to use?
      choices:
        - choice: Echo
          dependencies:
            - github.com/labstack/echo/v4
        - choice: Gorilla Mux
          dependencies:
            - github.com/gorilla/mux
        - choice: Gin
          dependencies:
            - github.com/gin-gonic/gin
```

</code-block>
</code-group>


### Multiple Select Question

If a question's `canSelectMultiple` is set to true, user can select more than one choice.

<code-group>
<code-block title="Terminal">

![Multiple Select Question](@images/multiple-select.png)

</code-block>

<code-block title="YAML">

```yaml
- name: Api Gateway
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: empty
  language: go
  questions:
    - direction: Select features you want
      canSelectMultiple: true
      choices:
        - choice: Elastic APM
          values:
            isElastic: true
        - choice: Swagger
          values:
            isSwagger: true
        - choice: Keycloak
          values:
            isKeycloak: true
```

</code-block>
</code-group>
