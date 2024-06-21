package auth

import (
	"chat-app/pkg/models"
	"chat-app/pkg/utils"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares a plain-text password with a hashed password
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// RegisterUser registers a new user with a username and hashedPassword
func RegisterUser(db *sql.DB, username, hashedPassword string) error {
	_, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		utils.Log.WithError(err).Error("Error inserting user into database")
		return err
	}

	return nil
}

// LoginUser logs in a user by verifying the username and password
func LoginUser(db *sql.DB, username, password string) (*models.Token, error) {
	user := &models.User{}

	row := db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	tokenString, err := GenerateJWT(user.ID, username)
	if err != nil {
		return nil, err
	}

	return &models.Token{Token: tokenString}, nil
}
