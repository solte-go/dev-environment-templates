.PHONY: run_watch_monitoring_tempo
run_watch:
	docker compose -f dc-prom-grafana-tempo.yaml watch

.PHONY: run_with_kafka_watch
run_with_kafka_watch:
	docker compose -f docker-compose-kafka.yaml -f dc-prom-grafana-tempo.yaml watch

.PHONY: run_with_kafka_watch
run_with_kafka_watch:
	docker compose -f docker-compose-kafka.yaml -f dc-prom-grafana-tempo.yaml watch
