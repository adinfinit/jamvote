package user

type ID string

type User struct {
	ID   ID
	Name string

	Email    string
	Facebook string
}
