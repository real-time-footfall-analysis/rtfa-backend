FROM golang:1.8.3-alpine

WORKDIR /go/src/github.com/real-time-footfall-analysis/rtfa-backend/

COPY . .

ARG RTFA_STATICDATA_DB_USER
ENV RTFA_STATICDATA_DB_USER=$RTFA_STATICDATA_DB_USER 
ARG RTFA_STATICDATA_DB_PASSWORD
ENV RTFA_STATICDATA_DB_PASSWORD=$RTFA_STATICDATA_DB_PASSWORD

RUN apk add --no-cache git
RUN go get -u github.com/go-pg/pg
RUN go get -u github.com/gorilla/mux

RUN go test -v ./...

RUN go build -o ~/go/bin/main .

ENTRYPOINT ~/go/bin/main

EXPOSE 80
