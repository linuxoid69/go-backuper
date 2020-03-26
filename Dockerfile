FROM alpine:3.9.5

COPY go-backuper /opt/

RUN chmod +x /opt/go-backuper && \
    apk update && \
    apk add postgresql-client tzdata && \
    mkdir -p /mnt/backups && \
    mkdir -p /mnt/data

CMD ["/opt/go-backuper"]
