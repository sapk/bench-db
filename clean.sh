#/bin/sh

docker-compose -p cassandra -f docker/cassandra/docker-compose.yml down
docker-compose -p mysql -f docker/mysql/docker-compose.yml down
docker-compose -p mariadb -f docker/mariadb/docker-compose.yml down
docker-compose -p mssql -f docker/mssql/docker-compose.yml down
#docker system prune -af