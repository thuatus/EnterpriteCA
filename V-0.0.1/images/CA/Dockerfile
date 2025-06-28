# Image form alpine
FROM alpine:latest

# Update SO and install required packages
RUN apk update && apk add --no-cache openssl openssh audit go

# Database installation
RUN apk add --no-cache mariadb mariadb-client bash

# Create directories for database, CA and application

RUN mkdir -p /run/mysqld /var/lib/mysql && \
    chown -R mysql:mysql /run/mysqld /var/lib/mysql


RUN mkdir /srv/CA /srv/ssl /srv/app

#Generates app certificate and key to guarantee ssl encryption
RUN openssl req -x509 -newkey rsa:4096 -keyout /srv/ssl/app.key -out /srv/ssl/app.crt -days 365 -nodes -subj "/CN=enterprise.app.com"

# Copy the application files to the container folder
COPY scriptCreateStructure.sql /srv/app/scriptCreateStructure.sql

# Initialize the database
RUN chown -R mysql:mysql /srv/app/scriptCreateStructure.sql
RUN chmod 644 /srv/app/scriptCreateStructure.sql
RUN mysql_install_db --user=mysql --basedir=/usr --datadir=/var/lib/mysql

# copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh



# Defines o entrypoint
ENTRYPOINT ["/entrypoint.sh"]



#Run sql script to create the database structure
COPY scriptCreateStructure.sql /docker-entrypoint-initdb.d/



# Set the working directory
WORKDIR /srv/app



