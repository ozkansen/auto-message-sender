package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"auto-message-sender/internal/models"
)

type messageRepository interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]models.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID, sendingStatus string) error
}

var _ messageRepository = (*MessagePostgresqlRepository)(nil)

type MessagePostgresqlRepository struct {
	conn *pgx.Conn
}

func NewMessagePostgresqlRepository(conn *pgx.Conn) *MessagePostgresqlRepository {
	return &MessagePostgresqlRepository{
		conn: conn,
	}
}

func (r *MessagePostgresqlRepository) GetUnsentMessages(ctx context.Context, limit int) ([]models.Message, error) {
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	messages, err := r.getUnsentMessages(ctx, tx, limit)
	if err != nil {
		return nil, err
	}
	err = r.setMessageStatusToPending(ctx, tx, messages)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessagePostgresqlRepository) getUnsentMessages(ctx context.Context, tx pgx.Tx, limit int) ([]models.Message, error) {
	rows, err := tx.Query(ctx, "SELECT message_id, phone_number, message_content, updated_at, created_at FROM messages WHERE sending_status = 'waiting' LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err2 := rows.Scan(
			&msg.MessageID,
			&msg.PhoneNumber,
			&msg.MessageContent,
			&msg.UpdatedAt,
			&msg.CreatedAt,
		)
		msg.SendingStatus = "waiting"
		if err2 != nil {
			return nil, err2
		}
		messages = append(messages, msg)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return messages, nil
}

func (r *MessagePostgresqlRepository) setMessageStatusToPending(ctx context.Context, tx pgx.Tx, messages []models.Message) error {
	for _, message := range messages {
		_, err := tx.Exec(ctx, "UPDATE messages SET sending_status = 'pending' WHERE message_id = $1", message.MessageID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MessagePostgresqlRepository) UpdateMessageStatus(ctx context.Context, messageID, sendingStatus string) error {
	_, err := r.conn.Exec(ctx, "UPDATE messages SET sending_status = $1 WHERE message_id = $2", sendingStatus, messageID)
	if err != nil {
		return err
	}
	return nil
}
