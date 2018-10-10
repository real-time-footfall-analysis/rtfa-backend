FROM golang:1.11.1-alpine

COPY . /go/src/github.com/real-time-footfall-analysis/rtfa-backend

RUN go test -v ./...

RUN go build -o ./bin/main /go/src/github.com/real-time-footfall-analysis/rtfa-backend/

ENTRYPOINT ./bin/main

EXPOSE 80
