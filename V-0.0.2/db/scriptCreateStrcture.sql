-- Create Database
CREATE DATABASE ca;


USE ca;

-- Criar Tables
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(45) NOT NULL,
    passwd CHAR(60) NOT NULL, -- store hashed passwords
    active BOOLEAN DEFAULT TRUE,
    mfaEnabled tinyint(1) DEFAULT '0',
    mfaUserKey varchar(60) DEFAULT NULL
);

CREATE TABLE ca (
    id INT AUTO_INCREMENT PRIMARY KEY,
    issuer VARCHAR(512),
    subject VARCHAR(512),
    public_Key TEXT,
    created DATETIME,
    expiration DATETIME,
    is_valid TINYINT(1) DEFAULT 1
);

