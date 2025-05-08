DB=user-service-db

CONTAINER_NAME=web-db-user-1

psql:
	docker exec -it $(CONTAINER_NAME) psql -U postgres -d $(DB)