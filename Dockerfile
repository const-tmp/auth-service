# syntax=docker/dockerfile:1

FROM golang:1.18 as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/auth cmd/auth/main.go
RUN go build -o bin/mgmt cmd/mgmt/main.go

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/bin/auth auth
COPY --from=build /app/bin/mgmt mgmt

USER nonroot:nonroot
