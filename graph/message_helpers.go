package graph

import (
	"strconv"
	"time"

	"github.com/samstringzz/alutamarket-backend/graph/model"
	"github.com/samstringzz/alutamarket-backend/internals/messages"
)

func convertMessageToGraphQL(m *messages.Message) *model.Message {
	if m == nil {
		return nil
	}

	// Convert ID to string
	id := strconv.FormatUint(uint64(m.ID), 10)
	sender := strconv.FormatUint(uint64(m.Sender), 10)
	chatID := strconv.FormatUint(uint64(m.ChatID), 10)

	// Format timestamps
	createdAt := m.CreatedAt.Format(time.RFC3339)
	var updatedAt *string
	if !m.UpdatedAt.IsZero() {
		t := m.UpdatedAt.Format(time.RFC3339)
		updatedAt = &t
	}

	// Convert media to string pointer
	var media *string
	if m.Media != nil {
		mediaStr := string(*m.Media)
		media = &mediaStr
	}

	return &model.Message{
		ID:        id,
		ChatID:    chatID,
		Content:   m.Content,
		Sender:    sender,
		Media:     media,
		IsRead:    m.IsRead,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
