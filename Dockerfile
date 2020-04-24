FROM golang:alpine
ADD . /go/src/cash
WORKDIR /go/src/cash
RUN go build -o cash
CMD ./cash
