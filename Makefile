db_prune:
	docker-compose -f=docker/docker-compose.yaml down -v
	docker-compose -f=docker/docker-compose.yaml up -d

run:
	go run ./... -d="host=localhost port=5432 user=myuser password=mypass dbname=mydatabase sslmode=disable"