# Build stage 
#- img version uptaded to use Go 1.25.3
FROM golang:1.25.3 AS builder 
WORKDIR /build
COPY . ./V-0.0.2/
WORKDIR /build/V-0.0.2/app
RUN go build -o /app/enterpriteca

# Final image  updated to use Alpine 3.22.2
FROM alpine:3.22.2

RUN apk add --no-cache ca-certificates openssl
WORKDIR /srv/app

# Copy built Go binary
COPY --from=builder /app/enterpriteca .

# Copy static, templates, and other needed folders
COPY app/static ./static
COPY app/templates ./templates
COPY app/index.html ./index.html
COPY ca /srv/ca
COPY db /srv/db

# (Optional) Generate SSL certs if needed
RUN mkdir -p /srv/ssl && \
    openssl req -x509 -newkey rsa:4096 -keyout /srv/ssl/app.key -out /srv/ssl/app.crt -days 365 -nodes -subj "/CN=enterprise.app.com"

EXPOSE 8443

CMD ["./enterpriteca"]


