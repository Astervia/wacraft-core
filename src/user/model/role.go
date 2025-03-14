package user_model

type Role string

var (
	Admin      Role = "admin"
	User       Role = "user"
	Automation Role = "automation"
	Developer  Role = "developer"
)
