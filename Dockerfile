FROM golang:1.14.3-alpine

RUN apk update && apk upgrade && apk add --no-cache git openssh
RUN go get github.com/cespare/reflex

WORKDIR /app

EXPOSE 8080

CMD reflex -s -- sh -c 'go run kernel.go'