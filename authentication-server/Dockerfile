FROM golang:1.11 as builder
WORKDIR /authentication-service
COPY . .
ENV GO111MODULE="on"
RUN go mod tidy
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go ./server.go

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /authentication-service/main /
COPY ./templates ./templates
COPY ./public ./public
COPY ./key.pem ./key.pem
COPY ./cert.pem ./cert.pem

EXPOSE 3000
CMD ["/main"]