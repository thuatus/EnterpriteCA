-- Create Database
CREATE DATABASE ca;


USE ca;

-- Criar Tables
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(45) NOT NULL,
    passwd CHAR(60) NOT NULL, -- store hashed passwords
    active BOOLEAN DEFAULT TRUE
);


