FROM docker.io/golang
LABEL maintainer="Ben Swinney"

COPY main.go /go/src/fcos-demo-server/

WORKDIR /go/src/fcos-demo-server

RUN go get -d -v ./... && \
    go install -v ./...

EXPOSE 8080/tcp

CMD ["fcos-demo-server"]