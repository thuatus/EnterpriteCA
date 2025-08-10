#\bin\bash
systemctl start mysqld.service
sleep 5
systemctl status mysqld.service

export MYSQL_ROOT_PASSWORD="passwd@G0" JWT_SECRET="1322190ufdjnfshupqwesad"

printenv | grep MYSQL
printenv | grep JWT

