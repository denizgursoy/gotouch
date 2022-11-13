# File

File content can be fetched from internet or provided in the YAML. A file must have `pathFromRoot` field to which Gotouch 
will create the file.

## Files from internet
```yaml
questions: #optional
  - direction: Do you want Dockerfile? #mandatory
    canSkip: true #if true, there must be at least one choice. 
    choices: #mandatory
      - choice: Yes
        files:
          - url: https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/Dockerfile
            pathFromRoot: Dockerfile #mandatory
```

## Files from content

```yaml
questions: #optional
  - direction: Do you want Dockerfile? #mandatory
    canSkip: true #if true, there must be at least one choice. 
    choices: #mandatory
      - choice: Yes
        files:
          - content: |
                FROM golang:1.12.0-alpine3.9
                RUN mkdir /app
                COPY . /app
                WORKDIR /app
                RUN go build -o main .
                CMD ["/app/main"]
            pathFromRoot: Dockerfile #mandatory
```
