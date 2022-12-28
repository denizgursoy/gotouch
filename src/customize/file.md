# File

File content can be fetched from internet or provided in the YAML. A file must have `pathFromRoot` field to which Gotouch 
will create the file.

## Files from Internet
```yaml
questions: #optional
  - direction: Do you want Dockerfile? #mandatory
    canSkip: true #if true, there must be at least one choice. 
    choices: #mandatory
      - choice: Yes # mandatory
        files:
          - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile # mandatory
            pathFromRoot: Dockerfile #mandatory
```

## Files from Content

```yaml
questions: #optional
  - direction: Do you want Dockerfile? #mandatory
    canSkip: true #if true, there must be at least one choice. 
    choices: #mandatory
      - choice: Yes # mandatory
        files: # content is mandatory
          - content: | 
                FROM golang:1.12.0-alpine3.9
                RUN mkdir /app
                COPY . /app
                WORKDIR /app
                RUN go build -o main .
                CMD ["/app/main"]
            pathFromRoot: Dockerfile #mandatory
```
