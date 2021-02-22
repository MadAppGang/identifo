FROM golang:1.16-alpine3.13 as builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/madappgang/identifo
COPY . ./
RUN go mod download
RUN go build -o /identifo .

FROM alpine:3.13.2
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /identifo .
#todo remove
COPY server-config.yaml ./
COPY static ./static
COPY jwt ./jwt
ENTRYPOINT ["./identifo"]