FROM golang:1.16-alpine AS builder

WORKDIR /app
COPY ./ ./
RUN go mod download


RUN go build -o /testapp

FROM gcr.io/distroless/base-debian10 AS runner

WORKDIR /

COPY --from=builder /testapp /testapp

USER nonroot:nonroot

ENTRYPOINT ["/testapp"]