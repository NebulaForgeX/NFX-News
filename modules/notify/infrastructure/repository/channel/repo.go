package channel

import (
	"context"
	"gorm.io/gorm"
	"nfxnews/modules/notify/domain/channel"
	"nfxnews/modules/notify/infrastructure/rdb/models"
)

type h struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *channel.Repo { return &channel.Repo{Create: &h{db: db}} }
func (x *h) New(ctx context.Context, c *channel.Channel) error {
	st := c.State()
	return x.db.WithContext(ctx).Create(&models.Channel{ID: st.ID, Kind: st.Kind, Name: st.Name, Enabled: st.Enabled, Config: st.Config, CreatedAt: st.CreatedAt, UpdatedAt: st.UpdatedAt}).Error
}
