FROM golang:1.13

RUN apt-get update 
RUN apt-get -y install redis-server

WORKDIR /go/src/server
COPY . .

RUN go get
RUN go build

CMD redis-server & ./main