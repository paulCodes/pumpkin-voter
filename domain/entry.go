package domain

type Entry struct {
	ID int64
	Title string
}

type EntryStore interface {
	// ListEntrys returns a list of Entrys, ordered by title.
	ListEntrys() ([]*Entry, error)

	// GetEntry retrieves a Entry by its ID.
	GetEntry(id int64) (*Entry, error)

	// AddEntry saves a given Entry, assigning it a new ID.
	AddEntry(b *Entry) (id int64, err error)

	// DeleteEntry removes a given Entry by its ID.
	DeleteEntry(id int64) error

	// UpdateEntry updates the entry for a given Entry.
	UpdateEntry(b *Entry) error

	// Close closes the database, freeing up any available resources.
	// TODO(cbro): Close() should return an error.
	Close()
}