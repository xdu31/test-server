FROM alpine:3.13

WORKDIR /bin

ENV PKG postgresql

RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache $PKG

COPY bin/server .
COPY pkg/server/pb/service.swagger.json www/swagger.json

ENTRYPOINT ["server", "--gateway.swaggerFile", "www/swagger.json"]
