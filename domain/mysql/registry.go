package mysql

import "gopkg.in/gorp.v1"

type Registry struct {
	Category MysqlCategoryStore
	Contest  MysqlContestStore
	Entry    MysqlEntryStore
	Vote     MysqlVoteStore
}

func NewRegistry(db *gorp.DbMap) Registry {
	store := NewMysqlStore(db)
	return Registry{
		Category: store.CategoryStore(),
		Contest:  store.ContestStore(),
		Entry:    store.EntryStore(),
		Vote:     store.VoteStore(),
	}
}
