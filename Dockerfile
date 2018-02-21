FROM golang:alpine

# ENV GOPATH /go
# ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" && \
        mkdir -p "$GOPATH/github.com/mikelsr"

COPY . "$GOPATH/src/github.com/mikelsr/sysmon"
RUN go install github.com/mikelsr/sysmon

CMD ["sysmon"]
