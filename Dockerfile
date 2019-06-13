FROM golang:1.10 as builder
# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/madappgang/identifo
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app .
COPY *config.yaml ./
COPY jwt/*.pem ./jwt/
COPY web/static ./static
COPY cmd/import/apps.json ./apps.json
COPY cmd/import/users.json ./users.json
COPY email_templates ./email_templates

CMD ["./app"]