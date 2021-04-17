package model

// Used to store a row of data from the membership table
type Membership struct {
	Id      uint64
	GroupId uint64
	UserId  uint64
}
