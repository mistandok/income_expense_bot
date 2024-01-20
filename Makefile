downapp:
	docker-compose -f docker-compose.yml down -v

startapp:
	docker-compose -f docker-compose.yml up -d --build

startservice:
	docker-compose up -d $(service_name)

migratecreate:
	migrate create -ext sql -dir db/migration -seq $(name)

.PHONY: downapp startapp migratecreate startservice