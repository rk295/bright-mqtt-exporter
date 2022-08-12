FROM golang:1.17-bullseye as build

WORKDIR /go/src/github.com/rk295/bright-mqtt-exporter

#ADD go.mod /go/src/github.com/rk295/bright-mqtt-exporter
#ADD go.sum /go/src/github.com/rk295/bright-mqtt-exporter
#RUN go get -d -v ./...

ADD . /go/src/github.com/rk295/bright-mqtt-exporter

RUN cd cmd/bright-mqtt-exporter && \
    go build -o /go/bin/bright-mqtt-exporter

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/bright-mqtt-exporter /
CMD ["/bright-mqtt-exporter"]
