package domain

type VoteStore interface {
	Add(course *Vote) error
	Replace(course *Vote) error
	Delete(id string) error
	FindById(id string) (Vote, error)
}


type Vote struct {
	ID string
	Category Category
	Entry Entry
	Contest Contest
}