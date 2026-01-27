package user

const (
	RoleUser    = "user"
	RolePremium = "premium"
	RoleAdmin   = "admin"
)

func (u *User) IsPremium() bool {
	return u.Role == RolePremium || u.Role == RoleAdmin
}
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
