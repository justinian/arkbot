FROM golang:1.15

ADD . /go/src/github.com/justinian/arkbot
WORKDIR /go/src/github.com/justinian/arkbot
RUN CGO_ENABLED=0 go build -a -o arkbot .


FROM alpine:latest
MAINTAINER Justin C. Miller <justin@devjustinian.com>

RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/justinian/arkbot/arkbot /
ENTRYPOINT ["/arkbot"]
