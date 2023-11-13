FROM golang:latest

RUN apt-get update && \
    mkdir -p /app

ADD ./bin/Mandatory4 /app/Mandatory4

RUN chmod +x /app/Mandatory4

EXPOSE 10000
ENTRYPOINT ["/app/Mandatory4"]