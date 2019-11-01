FROM golang:latest
COPY . /go/src/joy

WORKDIR /go/src/joy
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build
EXPOSE 8080
CMD ["./joy"]