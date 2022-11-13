# Question

Question allows your users to customize their projects. A question must have a direction and choices. Gotouch prompts
user to make a [choice](/customize/choice) for every question in the selected Project Structure.


## Yes/No question
If a question has only one choice and `canSkip` is true, it is evaluated as Yes/No question

```yaml
questions: #optional
  - direction: Do you want Dockerfile? #mandatory
    canSkip: true #if true, there must be at least one choice. 
    choices: #mandatory
      - choice: Yes
        files:
          - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
            pathFromRoot: Dockerfile
```
Will be displayed like:

![Yes/No Question](@images/yes-no-question.png)

## Multiple choice question

If a question has more than one choice it is evaluated as Multiple choice question

```yaml
- name: Api Gateway
  url: https://raw.githubusercontent.com/api/app/main/api-gateway.tar.gz # can be a tar.gz archive file
  language: go # must be go or golang for go projects, otherwise omit the field
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
  url: https://raw.githubusercontent.com/api/app/main/api-gateway.tar.gz # can be a tar.gz archive file
  language: go # must be go or golang for go projects, otherwise omit the field
  canSkip: true #if true, there must be at least one choice. 
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