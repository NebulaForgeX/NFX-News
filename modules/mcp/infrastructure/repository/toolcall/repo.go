package toolcall

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/mcp/domain/toolcall"
	"nfxnews/modules/mcp/infrastructure/rdb/models"
)

type h struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *toolcall.Repo { return &toolcall.Repo{Create: &h{db: db}} }
func (x *h) New(ctx context.Context, t *toolcall.ToolCall) error {
	st := t.State()
	return x.db.WithContext(ctx).Create(&models.ToolCall{ID: st.ID, ToolName: st.ToolName, Arguments: st.Arguments, OK: st.OK, ErrorMessage: st.ErrorMessage}).Error
}
