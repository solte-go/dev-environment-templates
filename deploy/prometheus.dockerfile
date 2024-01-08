FROM prom/prometheus:v2.48.1

COPY ./configs/prometheus/prometheus_cfg.yaml /etc/prometheus/cfg.yaml

ENTRYPOINT ["/bin/prometheus"]


