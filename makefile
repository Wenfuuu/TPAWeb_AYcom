DB=chat-service-db

CONTAINER_NAME=web-db-chat-1

psql:
	docker exec -it $(CONTAINER_NAME) psql -U postgres -d $(DB)