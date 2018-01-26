FROM golang:latest

COPY . /go/src/github.com/flameous/overlap-detection
WORKDIR /go/src/github.com/flameous/overlap-detection/cmd

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

VOLUME /host_dir/

EXPOSE 8080

CMD ["go-wrapper", "run"]
