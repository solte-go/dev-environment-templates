.PHONY: run_watch
run_watch:
	docker compose -f docker-compose.yaml -f docker-compose-overload.yaml watch

.PHONY: run_with_kafka_watch
run_with_kafka_watch:
	docker compose -f docker-compose.yaml -f docker-compose-kafka.yaml -f docker-compose-overload.yaml watch
