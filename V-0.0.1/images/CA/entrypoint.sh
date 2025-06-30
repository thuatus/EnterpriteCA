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

echo "🔧 Starting MariaDB server in backgound..."
mysqld --user=mysql --datadir=/var/lib/mysql --skip-networking &
pid="$!"

sleep 30
# Check if the database is already initialized and setup
for f in /docker-entrypoint-initdb.d/*.sql; do
    if [ -f "$f" ]; then
        echo "🔧 Importing SQL schema $f..."
        mysql -u root -p="${ROOT_PASS}" <"$f"

    fi
done

sleep 5
# Stop the temporary MariaDB server
mysqladmin -u root -p"${ROOT_PASS}" shutdown

# Start MariaDB normally
mysqld --user=mysql --datadir=/var/lib/mysql

# Start Go application

echo "Running go application...  " 
exec /srv/app/enterpriteca
