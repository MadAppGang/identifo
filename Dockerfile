FROM node:16 as node_builder

WORKDIR /identifo
COPY web_apps_src ./web_apps_src
COPY static ./static
RUN web_apps_src/update-admin.sh
RUN web_apps_src/update-web.sh


FROM golang:1.16-alpine3.13 as builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/madappgang/identifo
COPY . ./
RUN go mod download
RUN go build -o /identifo .

FROM alpine:3.13.2
RUN apk --no-cache add ca-certificates

WORKDIR /
COPY --from=builder /identifo .
COPY --from=node_builder /identifo/static ./static
COPY cmd/config-boltdb.yaml ./server-config.yaml
COPY jwt/test_artifacts/private.pem ./jwt/test_artifacts/private.pem
RUN mkdir -p /data

EXPOSE 8081/tcp

ENV IDENTIFO_ADMIN_LOGIN=admin@admin.com
ENV IDENTIFO_ADMIN_PASSWORD=password


ENTRYPOINT ["./identifo"]