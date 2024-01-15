# income_expense_bot

```app_start
docker-compose -f docker-compose.yml down -v
docker-compose -f docker-compose.yml up -d --build
```

```migrations
migrate create -ext sql -dir db/migration -seq init_schema
migrate -path db/migration -database "postgresql://username:password@localhost:5432/dbname?sslmode=disable" -verbose up
migrate -path db/migration -database "postgresql://username:password@localhost:5432/dbname?sslmode=disable" -verbose down
```
