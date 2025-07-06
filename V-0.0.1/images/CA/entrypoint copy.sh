#!/bin/sh
set -e

# Default values (can be overridden with environment variables)
DB_NAME=${MYSQL_DATABASE:-mydatabase}
DB_USER=${MYSQL_USER:-user}
DB_PASS=${MYSQL_PASSWORD:-secret123}
ROOT_PASS=${MYSQL_ROOT_PASSWORD:-root123}

# check if database is ready
if [! -d /var/lib/mysql/mysql ]; then
    echo "🔧 Initializing database directory..."
    mkdir -p /var/lib/mysql
    chown -R mysql:mysql /var/lib/mysql
    chmod 755 /var/lib/mysql
    mysql_install_db --user=mysql --datadir=/var/lib/mysql
    echo "Database directory initialized."
else
    echo "Database directory already initialized."
fi
# temporarily start MariaDB to set up root password and create database/user
echo "🚀 Starting temporary MariaDB server to set up root password and create database/user..."
mysqld_safe --skip-networking &
pid="$!"
# Wait for MariaDB to start
until mysqladmin ping --silent; do
    echo "⏳ Waiting for MariaDB to start..."
    sleep 1
done

# Setup root password, database, and user
echo "🔧 Setting up database and user..."
mysqld -u root <<-EOSQL
    
    ALTER USER 'root'@'localhost' IDENTIFIED BY '${ROOT_PASS}';
    CREATE DATABASE IF NOT EXISTS \`${DB_NAME}\`;
    CREATE USER IF NOT EXISTS '${DB_USER}'@'localhost' IDENTIFIED BY '${DB_PASS}';
    GRANT ALL PRIVILEGES ON \`${DB_NAME}\`.* TO '${DB_USER}'@'localhost';
    FLUSH PRIVILEGES;
EOSQL

echo "applyng app db congfiguration... "

# Import any .sql files
for f in /docker-entrypoint-initdb.d/*.sql; do
    if [ -f "$f" ]; then
        echo "🔧 Importing SQL schema $f..."
        mysql -u root -p"${ROOT_PASS}" <"$f"
    fi
done

# Stop the temporary MariaDB server
echo "🛑 Shutting down temporary MariaDB server..."
mysqladmin -u root -p"${ROOT_PASS}" shutdown

# Start MariaDB normally in the background
echo "🚀 Starting MariaDB in normal mode..."
mysqld --user=mysql --datadir=/var/lib/mysql

# Wait for MariaDB to be ready again
for i in {30..0}; do
    if mysqladmin -u root -p"${ROOT_PASS}" ping --silent; then
        break
    fi
    echo "⏳ Waiting for MariaDB to be ready (normal mode)..."
    sleep 1
done

if ! mysqladmin -u root -p"${ROOT_PASS}" ping --silent; then
    echo "❌ MariaDB did not start in normal mode."
    #exit 1
fi

# Start Go application (replace with exec if you want it to be PID 1)
echo "Running go application..."
exec /srv/app/enterpriteca
