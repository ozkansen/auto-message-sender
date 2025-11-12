# Automatic Message Sending System

## Description

Sample application that sends messages at periodic intervals

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
