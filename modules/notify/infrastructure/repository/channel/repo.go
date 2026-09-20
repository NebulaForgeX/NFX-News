package channel

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/notify/domain/channel"
	"nfxnews/modules/notify/infrastructure/rdb/models"
)

type handler struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *channel.Repo {
	impl := &handler{db: db}
	return &channel.Repo{Create: impl, Update: impl}
}

func toM(c *channel.Channel) *models.Channel {
	st := c.State()
	return &models.Channel{ID: st.ID, AccountID: st.AccountID, ProfileID: st.ProfileID, Kind: st.Kind, Name: st.Name, Enabled: st.Enabled, Config: st.Config, CreatedAt: st.CreatedAt, UpdatedAt: st.UpdatedAt}
}

func (handler *handler) New(ctx context.Context, c *channel.Channel) error {
	return handler.db.WithContext(ctx).Create(toM(c)).Error
}

func (handler *handler) Generic(ctx context.Context, c *channel.Channel) error {
	return handler.db.WithContext(ctx).Save(toM(c)).Error
}
