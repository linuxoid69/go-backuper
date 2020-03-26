FROM alpine:3.9.5

COPY go-backuper /opt/

RUN apk update && \
    apk add postgresql-client && \
    mkdir -p /mnt/backups && \
    mkdir -p /mnt/data

CMD ["/opt/go-backuper"]
