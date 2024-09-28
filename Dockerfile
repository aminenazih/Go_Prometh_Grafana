FROM golang:1.20

RUN apt-get update && apt-get install -y curl
RUN curl -L https://github.com/kyleconroy/sqlc/releases/download/v1.15.0/sqlc_1.15.0_linux_amd64.tar.gz | tar -C /usr/local/bin -xz

WORKDIR /app

CMD ["sqlc", "generate"]
