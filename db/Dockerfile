# package conversion scripts into the conversion container
FROM        infoblox/migrate:latest
RUN         /bin/mkdir -p /var/db
COPY        db/migrations /var/db/migrations
ENTRYPOINT ["/migrate", "--source", "file:///var/db/migrations"]
