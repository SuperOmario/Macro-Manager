package models

type User struct {
	userID int
	fName  string
	lName  string
	email  string
}

//Constructor for User struct
func NewUser(userID int, fName string, lName string, email string) User {
	return User{userID: userID, fName: fName, lName: lName, email: email}
}

//Getter for userID
func (u User) GetUserID() int {
	return u.userID
}
