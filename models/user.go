package models

type User struct {
	ID       int64
	Name     string
	Nickname string
	Password string
}

func (u *User) GetMyName() string {
	return u.Nickname
}
