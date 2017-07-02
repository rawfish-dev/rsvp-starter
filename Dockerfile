FROM golang:1.8

ADD . /go/src/github.com/rawfish-dev/rsvp-starter

WORKDIR /go/src/github.com/rawfish-dev/rsvp-starter

WORKDIR /go/src/github.com/rawfish-dev/rsvp-starter/server

RUN go install -a -v ./...

EXPOSE 6001

ENTRYPOINT ["/go/bin/server"]
