FROM golang:1.18.2-alpine3.16

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]