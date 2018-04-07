FROM golang:alpine

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" && \
        mkdir -p "$GOPATH/src/github.com/mikelsr"

COPY . "$GOPATH/src/github.com/mikelsr/sysmon"
RUN go install github.com/mikelsr/sysmon

CMD ["sysmon"]
