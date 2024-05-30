include .env
export

.PHONY: run_watch_monitoring_tempo
run_watch_monitoring_tempo:
	docker compose -f dc-prom-grafana-tempo.yaml watch

.PHONY: run_with_kafka_watch
run_with_kafka_watch:
	docker compose -f docker-compose-kafka.yaml -f dc-prom-grafana-tempo.yaml watch

.PHONY: run_with_kafka_sasl-scram
run_with_kafka_sasl:
	docker compose -f dc-kafka-sasl-scram.yaml up -d

.PHONY: run_with_kafka_ssl
run_with_kafka_ssl:
	docker compose -f dc-kafka-ssl.yaml up -d