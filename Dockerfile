FROM golang:1.20.4-alpine AS build

WORKDIR /build

COPY . .

RUN go mod tidy

RUN go build -o service

FROM alpine:latest

WORKDIR /app

COPY --from=build /build/service /app/service

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/service"]
