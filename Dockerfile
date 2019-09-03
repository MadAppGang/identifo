FROM golang:1.12.8 as builder

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/madappgang/identifo
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app .
COPY server-config.yaml ./
COPY jwt/*.pem ./jwt/
COPY web/static ./web/static
COPY cmd/import/apps.json ./apps.json
COPY cmd/import/users.json ./users.json
COPY email_templates ./email_templates
COPY admin_panel/build ./admin_panel/build

CMD ["./app"]