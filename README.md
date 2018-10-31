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

(cd docker/mariadb && docker-compose up -d)
MARIADB_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mariadb_database_1):3306)/" \
go test -bench=. -v
```
TODO regex bench

## All
```
docker-compose -p cassandra -f docker/cassandra/docker-compose.yml up -d && \
docker-compose -p mysql -f docker/mysql/docker-compose.yml up -d && \
docker-compose -p mariadb -f docker/mariadb/docker-compose.yml up -d
CASSANDRA_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cassandra_database_1) \
MYSQL_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql_database_1):3306)/" \
MARIADB_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mariadb_database_1):3306)/" \
go test -bench=. -v -benchtime=60s
```

## Clean
```
docker-compose -p cassandra -f docker/cassandra/docker-compose.yml down && \
docker-compose -p mysql -f docker/mysql/docker-compose.yml down && \
docker-compose -p mariadb -f docker/mariadb/docker-compose.yml down && \
docker system prune -af
```