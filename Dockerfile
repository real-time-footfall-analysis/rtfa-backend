FROM golang:1.11.1-alpine

COPY . /go/src/github.com/real-time-footfall-analysis/rtfa-backend

RUN go env

RUN ls -al /go

RUN ls -al ./src

RUN pwd

RUN go test -v ./...

RUN go build -o ./bin/main /go/src/github.com/real-time-footfall-analysis/rtfa-backend/


RUN ls -al ./
RUN ls -al ./bin

ENTRYPOINT ./bin/main

EXPOSE 80
