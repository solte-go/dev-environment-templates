FROM prom/prometheus:v2.43.0

COPY ./configs/prometheus/prometheus_cfg.yaml /etc/prometheus/cfg.yaml

ENTRYPOINT ["/bin/prometheus"]


