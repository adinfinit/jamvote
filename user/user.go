package user

type ID string

type User struct {
	ID   ID     `datastore:"KeyName"`
	Name string `datastore:",noindex"`

	Email    string `datastore:",noindex"`
	Facebook string `datastore:",noindex"`
}
