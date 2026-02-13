package billing_model

type Scope string

const (
	ScopeUser      Scope = "user"
	ScopeWorkspace Scope = "workspace"
)

func IsValidScope(s Scope) bool {
	return s == ScopeUser || s == ScopeWorkspace
}
