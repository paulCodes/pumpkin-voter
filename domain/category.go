package domain

type CategoryStore interface {
	Add(course *Category) error
	Replace(course *Category) error
	Delete(id string) error
	FindById(id string) (Category, error)
}


type Category struct {
	ID string
	Title string
	Active bool
}
