package factory

import (
	itemDomain "nfxnews/modules/news/domain/item"
	prefDomain "nfxnews/modules/news/domain/preference"
	itemRepo "nfxnews/modules/news/infrastructure/repository/item"
	prefRepo "nfxnews/modules/news/infrastructure/repository/preference"
	"nfxnews/pkgs/transaction"

	"gorm.io/gorm"
)

type TxRepoFactory struct{ db *gorm.DB }

func NewTxRepoFactory(db *gorm.DB) *TxRepoFactory { return &TxRepoFactory{db: db} }

func (f *TxRepoFactory) dbOr(uow transaction.UoW) *gorm.DB {
	if uow.DB != nil {
		return uow.DB
	}
	return f.db
}

func (f *TxRepoFactory) Item(uow transaction.UoW) *itemDomain.Repo {
	return itemRepo.NewRepo(f.dbOr(uow))
}

func (f *TxRepoFactory) Preference(uow transaction.UoW) *prefDomain.Repo {
	return prefRepo.NewRepo(f.dbOr(uow))
}
