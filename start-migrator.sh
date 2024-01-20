#!/bin/sh

wait_database()
{
  HOST=$1
  PORT=$2
  TYPE=$3

  echo "Waiting for $TYPE..."

  while ! nc -z $HOST $PORT; do
    sleep 0.1
  done

  echo "$TYPE started"
}

wait_database $DB_HOST $DB_PORT $DB_TYPE

migrate -path ./db/migration -database "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$DB_HOST:$DB_PORT/$POSTGRES_DB?sslmode=disable" -verbose up

echo "migrations for $DB_NAME was finished"