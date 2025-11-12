# Automatic Message Sending System

## Description

Sample application that sends messages at periodic intervals

*While developing the application, efforts were made to reduce external dependencies as much as possible.*

## Environment Variables

Postgresql example connection string for application database connection
- Name: "POSTGRESQL_DSN"
- Example value: "postgres://dbuser:dbpassword@postgresdb:5432/automessagesenderdb?sslmode=disable"

Redis example connection string for application cache connection
- Name: "REDIS_ADDR"
- Example value: "redis://localhost:6379/0"

Webhook.site example connection string for application webhook connection
- Name: "WEBHOOK_SITE_URL"
- Example value: "https://webhook.site/264d7ada-f7a7-40e9-8f30-eb0bde016436"

## How To Run

*Development default settings are available in docker-compose.yaml.
If you want to change them, you can change them from within the file.*

Require Docker & Docker Compose plugin

- Build and run
```bash
git clone https://github.com/ozkansen/auto-message-sender.git
cd auto-message-sender
docker-compose up --build
```

- Stop and remove
```bash
docker-compose down --remove-orphans --volumes
```

- Get Sent Messages
```bash
curl -X GET http://localhost:8080/messages | jq
```

Response:
```json
[
  {
    "message": "Accepted",
    "message_id": "3f846a61-2e99-42f9-a9ab-1e6cf1703476",
    "sent_at": "2025-11-12T01:09:51.133430722Z"
  },
  {
    "message": "Accepted",
    "message_id": "31a9f1f5-1ea2-4f74-bf8d-bcf4e587a482",
    "sent_at": "2025-11-12T01:09:51.229048492Z"
  }
]
```

- Start Auto Message Sender (When application started automatically starts)
```bash
curl -X POST http://localhost:8080/start
```

Response:
```
OK
```

- Stop Auto Message Sender
```bash
curl -X POST http://localhost:8080/stop
```

Response:
```
OK
```

- Health Check
```bash
curl -X GET http://localhost:8080/health
```

Response:
```
OK
```
