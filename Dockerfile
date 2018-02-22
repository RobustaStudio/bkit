FROM golang:alpine

RUN apk update && apk add git

RUN go get github.com/RobustaStudio/bkit

ENTRYPOINT ["bkit"]

WORKDIR /root/
