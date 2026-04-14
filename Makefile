db_up:
	docker-compose -f=docker/docker-compose.yaml up -d
	
db_prune:
	docker-compose -f=docker/docker-compose.yaml down -v
	db_up

run:
	go run ./... -d="host=localhost port=5432 user=myuser password=mypass dbname=mydatabase sslmode=disable"