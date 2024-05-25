#docker multistage build, done to avoid using docker-compose
FROM golang:1.22.3-alpine AS builder
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o cmd/main cmd/main.go

FROM mongo:latest 
EXPOSE 27017
COPY --from=builder /app/cmd/main /app/cmd/main
RUN mkdir /data/scripts
COPY servicesData.json versionsData.json startup.sh /data/scripts/
RUN chmod +x /data/scripts/startup.sh
ENTRYPOINT ["/data/scripts/startup.sh"]