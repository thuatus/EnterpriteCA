#!/bin/sh

# Default values (can be overridden with environment variables)
DB_NAME=${MYSQL_DATABASE:-mydatabase}
DB_USER=${MYSQL_USER:-user}
DB_PASS=${MYSQL_PASSWORD:-secret123}
ROOT_PASS=${MYSQL_ROOT_PASSWORD:-root123}

# Initialize the database if it doesn't exist
if [ ! -d "/var/lib/mysql/mysql" ]; then
    echo "🔧 Initializing database..."
    mysql_install_db --user=mysql --basedir=/usr --datadir=/var/lib/mysql

    # Start MariaDB temporarily for setup
    mysqld --user=mysql --bootstrap <<-EOSQL
        ALTER USER 'root'@'localhost' IDENTIFIED BY '${ROOT_PASS}';
        CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\`;
        CREATE USER IF NOT EXISTS '${DB_USER}'@'%' IDENTIFIED BY '${DB_PASS}';
        GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'%';
        FLUSH PRIVILEGES;
EOSQL
fi

# Start MariaDB normally
exec mysqld --user=mysql --datadir=/var/lib/mysql


# Routine do initialize and setup the database
if [ ! -d "/var/lib/mysql/mysql" ]; then
    echo "🔧 Inicializando banco de dados..."
    mysql_install_db --user=mysql --basedir=/usr --datadir=/var/lib/mysql

    # O MariaDB executará automaticamente os scripts em /docker-entrypoint-initdb.d
fi

