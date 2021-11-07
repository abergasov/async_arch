
stop:
	@echo "-- stop containers";
	docker container ls -q --filter name=arch* ; true
	@echo "-- drop containers"
	docker rm -f -v $(shell docker container ls -q --filter name=arch*) ; true

dev_up: stop
	@echo "RUN dev docker-compose.yml "
	docker-compose pull
	docker-compose up --build authDB taskDB broker zookeeper

migrate_new:
	migrate create -ext sql -dir migrations -seq data

migrate:
	migrate -path migrations/ -database "postgres://user:password@127.0.0.1:4666/arch_auth_db?query" up
