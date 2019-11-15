FROM golang:1.13

WORKDIR /go/src/client
COPY . .

RUN go get
RUN go build

CMD ./main