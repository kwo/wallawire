package model

import (
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	LoginTimeout = time.Hour * 24
)

// User defines a system user
type User struct {
	ID           string    `json:"id"`
	Disabled     bool      `json:"disabled"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"passwordHash"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

// MatchPassword checks if the given password matches the user password.
func (z *User) MatchPassword(password string) bool {
	hashedPassword, err := hex.DecodeString(z.PasswordHash)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
}

// SetPassword updates the password hashes.
func (z *User) SetPassword(password string) error {

	bhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	z.PasswordHash = hex.EncodeToString(bhash)

	return nil

}

// SessionToken is constructed from the JWT claims and stored in the request context
type SessionToken struct {
	SessionID string    `json:"sessionID"`
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Roles     []string  `json:"roles"`
	Issued    time.Time `json:"-"`
	Expires   time.Time `json:"-"`
}

// HasRole tests if the user has been assigned the given role
func (z *SessionToken) HasRole(role string) bool {
	if len(z.Roles) > 0 {
		for _, value := range z.Roles {
			if value == role {
				return true
			}
		}
	}
	return false
}

func ToSessionToken(sessionID string, u *User, roles []UserRole, issued, expires time.Time) *SessionToken {

	if u == nil {
		return nil
	}

	var rolenames []string
	for _, role := range roles {
		rolenames = append(rolenames, role.Name)
	}

	return &SessionToken{
		SessionID: sessionID,
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Roles:     rolenames,
		Issued:    issued,
		Expires:   expires,
	}

}
