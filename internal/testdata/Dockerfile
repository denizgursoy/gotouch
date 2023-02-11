§FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY ./ ./

RUN go mod download


RUN CGO_ENABLED=0 go build -o /testapp

FROM gcr.io/distroless/base-debian10 AS runner

WORKDIR /

COPY --from=builder /testapp /testapp

EXPOSE 8090

USER nonroot:nonroot

ENTRYPOINT ["/testapp"]