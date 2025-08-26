package user

type User struct {
	Username   string
	Password   string
	Email      string
	Enabled    bool
	Role       string // e.g., "admin" or "user"
	mfaSecret  string
	mfaEnabled bool
}
