FROM golang:1.12.1

RUN apt-get update && apt-get install make

WORKDIR $GOPATH/src/github.com/port-domain

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./vendor ./vendor
COPY Makefile ./main.go ./

RUN make && \
    cp ./port /usr/local/bin/ && \
    rm -rf /go/src/github.com

WORKDIR /usr/local/bin/

ENV BIND 0.0.0.0:8000

EXPOSE 8000

ENTRYPOINT ["port"]

