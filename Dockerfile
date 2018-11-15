FROM golang:1.8.3-alpine

WORKDIR /go/src/github.com/real-time-footfall-analysis/rtfa-backend/

COPY . .

RUN apk add --no-cache git

RUN go get -u -d github.com/go-pg/pg
RUN (cd ../../go-pg/pg/ && git checkout ae5d5e7df4b2e598390e10b66b849c6af94f092b)

RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/aws/aws-sdk-go/aws
RUN go get -u github.com/aws/aws-sdk-go/aws/session
RUN go get -u github.com/aws/aws-sdk-go/service/kinesis
RUN go get -u github.com/aws/aws-sdk-go/service/dynamodb

RUN go test -v ./...

RUN go build -o ~/go/bin/main .

ENTRYPOINT ~/go/bin/main

EXPOSE 80
