# Question

Question allows your users to customize their projects. A question must have a direction and choices. Gotouch prompts
user to make a [choice](./choice) for every question in the selected Project Structure.

**`direction`**: Direction is displayed for question. It is a mandatory field.

## Yes/No question
If a question has only one choice and `canSkip` is true, it is evaluated as Yes/No question

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
Will be displayed like:

![Yes/No Question](@images/yes-no-question.png)

## Multiple choice question

If a question has more than one choice it is evaluated as Multiple choice question

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

Will be displayed like:

![Multiple Choice Question](@images/multiple-choice.png)

## None of above question

If a question has more than one choice and `canSkip` is true, `None of above` option will be also added 
as a choice. 

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

Will be displayed like:


![None of Above Question](@images/none-of-above.png)


## Multiple Select Question

If a question's `canSelectMultiple` is set to true, user can select more than one choice.

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

Will be displayed like:


![Multiple Select Question](@images/multiple-select.png)