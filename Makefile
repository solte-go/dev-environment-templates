include .env
export

#### --- HELP --- ####
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


.PHONY: run_watch_monitoring_tempo
run_watch_monitoring_tempo:
	docker compose -f dc-base.yaml -f dc-grafana.yaml -f dc-tempo.yaml watch

.PHONY: run_with_kafka_watch
run_with_kafka_watch:
	docker compose -f docker-compose-kafka.yaml -f dc-prom-grafana-tempo.yaml watch

.PHONY: run_with_kafka_sasl-scram
run_with_kafka_sasl:
	docker compose -f dc-kafka-sasl-scram.yaml up -d

.PHONY: run_with_kafka_ssl
run_with_kafka_ssl:
	docker compose -f dc-kafka-ssl.yaml up -d

.PHONY: teardown_with_kafka_ssl
teardown_with_kafka_ssl:
	docker compose -f dc-kafka-ssl.yaml down -v