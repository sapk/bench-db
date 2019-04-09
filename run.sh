#/bin/sh

#Cassandra
docker-compose -p cassandra -f docker/cassandra/docker-compose.yml up -d
sleep 2m
CASSANDRA_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cassandra_database_1) \
go test -bench=. -v -benchtime=60s
CASSANDRA_IP=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' cassandra_database_1) \
go test -bench=. -v -benchtime=300s
docker-compose -p cassandra -f docker/cassandra/docker-compose.yml down

docker-compose -p mysql -f docker/mysql/docker-compose.yml up -d
sleep 2m
MYSQL_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql_database_1):3306)/" \
go test -bench=. -v -benchtime=60s
MYSQL_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql_database_1):3306)/" \
go test -bench=. -v -benchtime=300s
docker-compose -p mysql -f docker/mysql/docker-compose.yml down

docker-compose -p mariadb -f docker/mariadb/docker-compose.yml up -d
sleep 2m
MARIADB_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mariadb_database_1):3306)/" \
go test -bench=. -v -benchtime=60s
MARIADB_URL="root:password@tcp($(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mariadb_database_1):3306)/" \
go test -bench=. -v -benchtime=300s
docker-compose -p mariadb -f docker/mariadb/docker-compose.yml down

docker-compose -p mssql -f docker/mssql/docker-compose.yml up
sleep 2m
MSSQL_URL="sqlserver://sa:1Password@$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mssql_database_1):1433/" \
go test -bench=. -v -benchtime=60s
MSSQL_URL="sqlserver://sa:1Password@$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mssql_database_1):1433/" \
go test -bench=. -v -benchtime=300s
docker-compose -p mssql -f docker/mssql/docker-compose.yml down


