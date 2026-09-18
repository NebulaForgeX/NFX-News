package channel

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/notify/domain/channel"
	"nfxnews/modules/notify/infrastructure/rdb/models"
)

type handler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *channel.Repo { return &channel.Repo{Create: &handler{db: db}} }
func (handler *handler) New(ctx context.Context, c *channel.Channel) error {
	st := c.State()
	return handler.db.WithContext(ctx).Create(&models.Channel{ID: st.ID, AccountID: st.AccountID, ProfileID: st.ProfileID, Kind: st.Kind, Name: st.Name, Enabled: st.Enabled, Config: st.Config, CreatedAt: st.CreatedAt, UpdatedAt: st.UpdatedAt}).Error
}
