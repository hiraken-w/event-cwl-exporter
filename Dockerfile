FROM golang:latest as builder
WORKDIR /go/src/github.com/hiraken-w/event-cwl-exporter
COPY . .
# Set Environment Variable
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# Build
RUN go build -o app ./cmd

FROM amazonlinux:2 as amazonlinux
# Runtime Container
FROM scratch
COPY --from=amazonlinux /etc/ssl/certs/ca-bundle.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/hiraken-w/event-cwl-exporter/app app
ENTRYPOINT ["/app"]
