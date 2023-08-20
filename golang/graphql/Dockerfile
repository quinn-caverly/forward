FROM golang:1.20-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN apk update && \
    apk add --no-cache git

RUN go get ./...

RUN go build -o main .

CMD ["/app/main"]
