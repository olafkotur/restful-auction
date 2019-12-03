# Build the image
FROM golang:1.13-alpine3.10 as build
WORKDIR /go/src/app
COPY . .
RUN go build

# Copy over only binary
FROM alpine:3.10
COPY --from=build /go/src/app/main .
EXPOSE 8080
CMD ./main
