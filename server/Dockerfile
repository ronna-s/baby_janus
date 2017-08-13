FROM golang:latest
WORKDIR /go/src/github.com/wwgberlin/baby_janus/server
RUN cd /go/src/github.com/wwgberlin/baby_janus/server
CMD go test ./... && go build && /go/src/github.com/wwgberlin/baby_janus/server/server
EXPOSE 8080