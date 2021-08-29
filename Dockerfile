FROM golang:1.14.2-alpine

ADD . /go/src/github.com/dewski/awair_exporter
RUN go install github.com/dewski/awair_exporter

EXPOSE 8181

ENTRYPOINT /go/bin/awair_exporter
