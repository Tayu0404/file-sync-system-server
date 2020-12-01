FROM mariadb:10.5

COPY custom.cnf /etc/mysql/conf.d
COPY initialize.sql /docker-entrypoint-initdb.d/
COPY filetype.csv /docker-entrypoint-initdb.d/