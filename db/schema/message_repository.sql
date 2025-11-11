CREATE TYPE sending_status AS ENUM (
    'waiting',
    'pending',
    'sent',
    'failed'
    );

CREATE TABLE IF NOT EXISTS messages
(
    message_id      UUID PRIMARY KEY,
    phone_number    VARCHAR(20),
    message_content VARCHAR(160),
    sending_status  sending_status,
    created_at      TIMESTAMP,
    updated_at      TIMESTAMP
);
