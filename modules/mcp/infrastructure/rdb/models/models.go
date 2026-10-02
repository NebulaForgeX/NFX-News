package models

import (
	"github.com/google/uuid"
)

type ToolCall struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ToolName     string
	Arguments    []byte `gorm:"type:jsonb"`
	OK           bool
	ErrorMessage *string
}

func (ToolCall) TableName() string { return "mcp.tool_calls" }
