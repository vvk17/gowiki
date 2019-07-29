FROM golang:alpine3.10

WORKDIR /go/src/gowiki
COPY . .

EXPOSE 8080/tcp

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["gowiki"]