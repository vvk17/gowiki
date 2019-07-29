FROM golang:alpine3.10

WORKDIR /go/src/gowiki
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["gowiki"]