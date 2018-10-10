FROM golang:1.11.1-alpine

WORKDIR /go/src/github.com/real-time-footfall-analysis/rtfa-backend/

COPY . .

RUN go get -u github.com/gorilla/mux

RUN go test -v ./...

RUN go build -o ~/go/bin/main .

ENTRYPOINT ~/go/bin/main

EXPOSE 80
