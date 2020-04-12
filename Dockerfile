###################
##  build stage  ##
###################
FROM golang:1.13.0-alpine as builder
WORKDIR /redis-golang-kubernetes
COPY . .
RUN go build -v -o redis-golang-kubernetes

##################
##  exec stage  ##
##################
FROM alpine:3.10.2
WORKDIR /app
COPY --from=builder /redis-golang-kubernetes /app/
CMD ["./redis-golang-kubernetes"]
