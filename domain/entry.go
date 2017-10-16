package domain

type EntryStore interface{
	Add(course *Entry) error
	Replace(course *Entry) error
	Delete(id string) error
	FindById(id string) (Entry, error)
}

type Entry struct {
	ID string
	Title string
	Categories []Category
}