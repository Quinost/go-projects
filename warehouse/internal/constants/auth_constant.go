package constants

type contextKey string

const (
	Authorization string     = "Authorization"
	Bearer        string     = "Bearer "
	ClaimsKey     contextKey = "claims"

	ClaimSub      = "sub"
	ClaimExp      = "exp"
	ClaimUserName = "user_name"
)
