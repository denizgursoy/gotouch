FROM golang:alpine

RUN go install github.com/denizgursoy/gotouch/cmd/gotouch@latest

WORKDIR /out

ENTRYPOINT [ "gotouch" ]
