package domain

type Category struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	Active bool   `db:"active"`
}

type Contest struct {
	Id          string `db:"id"`
	Title       string `db:"title"`
	CategoryIds string `db:"category_ids"`
	Active      bool   `db:"active"`
}

type Entry struct {
	Id          string `db:"id"`
	Title       string `db:"title"`
	CategoryIds string `db:"category_ids"`
	ContestId   string `db:"contest_id"`
}

type Vote struct {
	Id         string `db:"id"`
	CategoryId string `db:"category_id"`
	EntryId    string `db:"entry_id"`
	ContestId  string `db:"contest_id"`
}
