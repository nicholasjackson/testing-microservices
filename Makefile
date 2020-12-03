run_db:
	docker run -d --name db -e POSTGRES_DB=root -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -p 5432:5432 shipyardrun/postgres:9.6

stop_db:
	docker rm -f db

restart_db: stop_db run_db
