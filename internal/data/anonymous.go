package data

var AnonymousUser = &User{}

func (u *User) IsAnon() bool {
	return u == AnonymousUser
}
