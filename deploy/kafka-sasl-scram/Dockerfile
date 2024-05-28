# It is recommened to use same version as kafka broker is used.
# So no additional images are pulled.
FROM confluentinc/cp-kafka:7.5.3

WORKDIR usr/bin

# Once it is executed, this container is not needed.
COPY startup.sh startup.sh
ENTRYPOINT ["./startup.sh"]