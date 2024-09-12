FROM prom/prometheus:v2.53.2

COPY ./configs/prometheus/prometheus_cfg.yaml /etc/prometheus/cfg.yaml

ENTRYPOINT ["/bin/prometheus"]


