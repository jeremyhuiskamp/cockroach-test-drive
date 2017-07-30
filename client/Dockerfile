FROM golang:alpine

COPY ./migrate /migrate

WORKDIR /go/src/app
COPY . .

RUN go-wrapper install

CMD ["go-wrapper", "run"]
