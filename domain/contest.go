package domain

type ContestStore interface {
	Add(course *Contest) error
	Replace(course *Contest) error
	Delete(id string) error
	FindById(id string) (Contest, error)
}



type Contest struct {
	ID string
	Title string
	Categories []Category
	Active bool
}