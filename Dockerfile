FROM golang:1.10-alpine AS builder
RUN apk add git
RUN go get -u github.com/golang/dep/cmd/dep
RUN mkdir -p $GOPATH/src/build/
COPY Gopkg.* $GOPATH/src/build/
WORKDIR $GOPATH/src/build

COPY main.go $GOPATH/src/build/
RUN dep ensure
RUN go build -o /bin/urlshorten

FROM alpine
COPY --from=builder /bin/urlshorten /bin/urlshorten
RUN mkdir -p /boltdb-data
EXPOSE 8080
ENTRYPOINT ./bin/urlshorten

