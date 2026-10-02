package factory

import (
	"gorm.io/gorm"
	keywordDomain "nfxnews/modules/report/domain/keyword"
	snapshotDomain "nfxnews/modules/report/domain/snapshot"
	keywordRepo "nfxnews/modules/report/infrastructure/repository/keyword"
	snapshotRepo "nfxnews/modules/report/infrastructure/repository/snapshot"
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
func (f *TxRepoFactory) Keyword(uow transaction.UoW) *keywordDomain.Repo {
	return keywordRepo.NewRepo(f.dbOr(uow))
}
func (f *TxRepoFactory) Snapshot(uow transaction.UoW) *snapshotDomain.Repo {
	return snapshotRepo.NewRepo(f.dbOr(uow))
}
