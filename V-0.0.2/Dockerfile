# Build stage 
#- img version uptaded to use Go 1.25.3
FROM golang:1.24.2 AS builder 
WORKDIR /build
COPY . .
WORKDIR /build/V-0.0.2/app

# static compilation

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64


RUN go build -o /app/enterpriteca
RUN ls -la /app/enterpriteca

# Final image  updated to use Alpine 3.22.2
FROM alpine:3.22.2

RUN apk add --no-cache ca-certificates openssl
WORKDIR /srv/app

# Copy built Go binary
COPY --from=builder /app/enterpriteca .

# Copy static, templates, and other needed folders
COPY V-0.0.2/app/static ./static
COPY V-0.0.2/app/templates ./templates
COPY V-0.0.2/app/index.html ./index.html
COPY V-0.0.2/ca /srv/ca
COPY V-0.0.2/db /srv/db

# (Optional) Generate SSL certs if needed
RUN mkdir -p /srv/ssl && \
    openssl req -x509 -newkey rsa:4096 -keyout /srv/ssl/app.key -out /srv/ssl/app.crt -days 365 -nodes -subj "/CN=enterprise.app.com"


RUN chmod +x ./enterpriteca

EXPOSE 8443

CMD ["./enterpriteca"]