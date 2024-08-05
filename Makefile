run: run-articles-service run-auth-service run-metrics

add-docker-network:
	docker network create articles-network

run-articles-service:
	docker-compose -f articles-service/docker-compose.yml up -d

run-auth-service:
	docker-compose -f auth-service/docker-compose.yml up -d

run-metrics:
	docker-compose -f collecting-metrics/docker-compose.yml up -d
