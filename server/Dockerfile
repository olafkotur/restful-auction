FROM golang:1.13

WORKDIR /go/src/server
COPY . .

RUN go build

EXPOSE 8080

CMD ./main
