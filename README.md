# EnterpriteCA

**EnterpriteCA – Enterprise Private Certificate Authority** is a Private Key Infrastructure (PKI) solution designed for secure certificate management in private environments. This project provides tools for creating, signing, and managing digital certificates, making it suitable for development and testing environments.

The purpose of this project is to learn Go and apply that knowledge in a real-world scenario related to the cybersecurity field.

## Version

**v0.0.1** – Initial version with the basic implementation of a Certificate Authority, a database, and a web application for development and learning purposes.

## Features

- Generate and maintain Certificate Authorities (CAs)
- Issue and revoke server certificates
- Support for X.509 certificates
- Secure storage of private keys
- Web interface for CA interaction

## Getting Started

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/EnterpriteCA.git
   cd EnterpriteCA

   Execute the image
```
   podman run -e MYSQL_DATABASE=ca \
           -e MYSQL_USER=admin \
           -e MYSQL_PASSWORD=123456 \
           -e MYSQL_ROOT_PASSWORD=passwd@G0 \
           eterpriteca
```