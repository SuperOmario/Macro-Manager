package models

type User struct {
	UserID int64
	FName  string
	LName  string
	Email  string
}

//Constructor for User struct
func NewUser(userID int64, fName string, lName string, email string) User {
	return User{UserID: userID, FName: fName, LName: lName, Email: email}
}

//Getter for userID
func (u User) GetUserID() int64 {
	return u.UserID
}
