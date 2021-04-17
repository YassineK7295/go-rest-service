package model

// Used to store one row of data from the user table
type User struct {
	Id        uint64
	FirstName string
	LastName  string
	UserId    string
}
