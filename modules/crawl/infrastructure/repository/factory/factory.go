package factory

import (
	"gorm.io/gorm"
	sessionDomain "nfxnews/modules/crawl/domain/session"
	sessionRepo "nfxnews/modules/crawl/infrastructure/repository/session"
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
func (f *TxRepoFactory) Session(uow transaction.UoW) *sessionDomain.Repo {
	return sessionRepo.NewRepo(f.dbOr(uow))
}
