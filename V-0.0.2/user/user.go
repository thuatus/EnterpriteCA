package user

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"image/png"
	"log"

	"github.com/pquerna/otp/totp"
	"github.com/thuatus/EnterpriteCA/V-0.0.2/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username   string
	Password   string
	Email      string
	Enabled    bool
	Role       string // e.g., "admin" or "user"
	MfaSecret  string
	MfaEnabled bool
	MfaQrImg   template.URL
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
		u.MfaEnabled = false
		return false
	}
	u.MfaEnabled = mfaEnabled
	return u.MfaEnabled
}

// Generates the MFA secret for the user
func (u *User) GenerateMFASecret() error {
	// Placeholder for generating MFA secret
	// In a real implementation, this would generate a secure random secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "EnterpriteCA",
		AccountName: u.Username,
	})
	if err != nil {
		log.Println("Error generating MFA secret:", err)
		return err
	}
	u.MfaSecret = key.Secret()
	img, _ := key.Image(200, 200)

	// Convert image to base64
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		log.Println("Error encoding image to PNG:", err)
		return err
	}
	encodedImg := base64.StdEncoding.EncodeToString(buf.Bytes())
	u.MfaQrImg = template.URL("data:image/png;base64," + encodedImg)
	log.Printf("MFA QR Code generated successfully : %s", u.MfaQrImg)

	return nil

}

// enables MFA for the user
func (u *User) EnableMFA(secret string) {
	u.MfaSecret = secret
	u.MfaEnabled = true

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection is not established")
		return
	}
	defer dbConn.Close()

	query := "UPDATE ca.users SET mfaUserKey = ?, mfaEnabled = ? WHERE name = ?"
	_, err = dbConn.Exec(query, u.MfaSecret, u.MfaEnabled, u.Username)
	if err != nil {
		log.Println("Error enabling MFA:", err)
		return
	}
	log.Println("MFA enabled successfully for user:", u.Username)

}

func (u *User) FirstValidateMFAToken(token string, secret string) bool {
	valid := totp.Validate(token, secret)
	if valid {
		u.EnableMFA(u.MfaSecret)
	}
	return valid
}
