FROM golang:alpine

WORKDIR /go/src/app
COPY ./src/client .

RUN go-wrapper install

CMD ["go-wrapper", "run"]
