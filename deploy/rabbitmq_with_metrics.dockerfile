FROM rabbitmq:3.9-management

RUN rabbitmq-plugins enable --offline rabbitmq_management rabbitmq_prometheus