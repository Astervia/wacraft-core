package user_model

type SearchableUserColumn string

var (
	Name       SearchableUserColumn = "name"
	Email      SearchableUserColumn = "email"
	RoleColumn SearchableUserColumn = "role"
)

func (t SearchableUserColumn) IsValid() bool {
	switch t {
	case Name, Email, RoleColumn:
		return true
	default:
		return false
	}
}
