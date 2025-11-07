-- Create Database
CREATE DATABASE IF NOT EXISTS ca;


USE ca;

-- Criar Tables
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(45) NOT NULL,
    passwd CHAR(60) NOT NULL, -- store hashed passwords
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS ca (
    id INT AUTO_INCREMENT PRIMARY KEY,
    issuer VARCHAR(512),
    subject VARCHAR(512),
    public_Key TEXT,
    created DATETIME,
    expiration DATETIME,
    is_valid TINYINT(1) DEFAULT 1
);

