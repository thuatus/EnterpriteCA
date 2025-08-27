package user

import (
	"log"

	"github.com/thuatus/EnterpriteCA/V-0.0.2/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username   string
	Password   string
	Email      string
	Enabled    bool
	Role       string // e.g., "admin" or "user"
	mfaSecret  string
	mfaEnabled bool
}

// verifies if the user exists in the database
func (u *User) Exists() bool {
	// Placeholder for database check
	// In a real implementation, this would query the database
	if db.UserExists(u.Username) {
		return true
	}

	return false
}

// saves the user to the database
func (u *User) SaveUser(password string) error {
	// Placeholder for database save operation
	// In a real implementation, this would insert or update the user in the database
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	// Create the admin user
	err = db.AddUser(u.Username, string(hash))
	if err != nil {
		log.Println("Error creating admin user:", err)
		return err

	}

	log.Println("Admin user created successfully")
	//http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

// Verifies if MFA is enabled for the user
func (u *User) IsMFAEnabled() bool {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection is not established")
		return false
	}
	defer dbConn.Close()
	var mfaEnabled bool
	query := "SELECT mfaEnabled FROM ca.users WHERE name = ?"
	err = dbConn.QueryRow(query, u.Username).Scan(&mfaEnabled)
	if err != nil {
		log.Println("Error checking MFA status:", err)
		u.mfaEnabled = false
		return false
	}
	u.mfaEnabled = mfaEnabled
	return u.mfaEnabled
}

// enables MFA for the user
func (u *User) EnableMFA(secret string) {
	u.mfaSecret = secret
	u.mfaEnabled = true
}
