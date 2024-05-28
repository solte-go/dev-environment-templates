#!/bin/bash

# Wait until the Kafka server is started
sleep 30

# Add ACL for user 'admin'
kafka-configs --bootstrap-server kafka:9092 --entity-type users --entity-name admin --alter --add-config 'SCRAM-SHA-512=[password=admin-secret]'

# Keep script in foreground
wait


