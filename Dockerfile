FROM golang:1.11.1-alpine

COPY . .

RUN go test -v ./...

RUN go build -o hello_world

ENTRYPOINT ./hello_world

EXPOSE 80
