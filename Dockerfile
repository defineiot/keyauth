
#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/github.com/defineiot/keyauth
COPY . .
RUN apk add --no-cache make
RUN apk add git
RUN make build
#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/defineiot/keyauth/keyauthd /keyauthd
ENTRYPOINT ./keyauthd
EXPOSE 8080
