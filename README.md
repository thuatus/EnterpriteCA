# EnterpriteCA

**EnterpriteCA – Enterprise Private Certificate Authority** is a Private Key Infrastructure (PKI) solution designed for secure certificate management in private environments. This project provides tools for creating, signing, and managing digital certificates, making it suitable for development and testing environments.

The purpose of this project is to learn Go and apply that knowledge in a real-world scenario related to the cybersecurity field.

## Versions

**v0.0.1** – Initial version with the basic implementation of a Certificate Authority, a database, and a web application for development and learning purposes.

**v0.0.2** - New features and bug fixes - acf0b9872407b46b66f146e0b181edb78d958924

## Features

- Generate and maintain Certificate Authorities (CAs)
- Issue and revoke server certificates
- Support for X.509 certificates
- Secure storage of private keys
- Web interface for CA interaction
- 2FA OTP Authentication

## Getting Started

1. Clone the repository and navigate to the images folder.
   ```bash
   git clone https://github.com/yourusername/EnterpriteCA.git
   cd EnterpriteCA/V-0.0.2/images/CA
   ```
2. Set the environment variables in `.env` file
   ```bash
   vim .env
   ```
3. Run docker compose
   ```bash
   docker compose up -d
   ```
4. Access https://localhost:8443 and configure your CA:

5. After that, you can issue, view, and revoke your certificates.
   