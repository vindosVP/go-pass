FROM golang:1.22.0

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy
RUN go build -o server ./cmd/server/main.go
RUN go build -o migrator ./cmd/migrator/main.go