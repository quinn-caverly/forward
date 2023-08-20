FROM golang:1.12.0-alpine3.9

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN apk update && \
    apk add --no-cache git

RUN go get ./...

RUN go build -o main .

CMD ["/app/main"]
