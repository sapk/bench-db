# bench-db-ripe-atlas
use ripe atlas data to benchmark various database

## Run specific bench
```
(cd docker/cassandra && docker-compose up -d)
CASSANDRA_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cassandra_database_1) \
go test -bench=. -v

(cd docker/mysql && docker-compose up -d)
MYSQL_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql_database_1):3306)/" \
go test -bench=. -v
```
TODO regex bench

## All
```
CASSANDRA_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cassandra_database_1) \
MYSQL_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql_database_1):3306)/" \
go test -bench=. -v -benchtime=60s
```