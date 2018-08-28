FROM golang:1.9-alpine

label MAINTAINER="Shantanu Deshpande <shantanud106@gmail.com>"

RUN apk add --no-cache vim git

RUN go get github.com/tsenart/vegeta

CMD [ "tail", "-F", "-n0", "/etc/hosts" ]
