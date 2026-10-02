package factory

import (
	"gorm.io/gorm"
	channelDomain "nfxnews/modules/notify/domain/channel"
	deliveryDomain "nfxnews/modules/notify/domain/delivery"
	channelRepo "nfxnews/modules/notify/infrastructure/repository/channel"
	deliveryRepo "nfxnews/modules/notify/infrastructure/repository/delivery"
	"nfxnews/pkgs/transaction"
)

type TxRepoFactory struct{ db *gorm.DB }

func NewTxRepoFactory(db *gorm.DB) *TxRepoFactory { return &TxRepoFactory{db: db} }
func (f *TxRepoFactory) dbOr(uow transaction.UoW) *gorm.DB {
	if uow.DB != nil {
		return uow.DB
	}
	return f.db
}
func (f *TxRepoFactory) Channel(uow transaction.UoW) *channelDomain.Repo {
	return channelRepo.NewRepo(f.dbOr(uow))
}
func (f *TxRepoFactory) Delivery(uow transaction.UoW) *deliveryDomain.Repo {
	return deliveryRepo.NewRepo(f.dbOr(uow))
}
