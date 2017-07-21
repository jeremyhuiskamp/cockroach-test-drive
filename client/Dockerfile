FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN go-wrapper install

CMD ["go-wrapper", "run"]
