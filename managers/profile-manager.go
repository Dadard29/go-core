package managers

import (
	"errors"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"time"
)

func ProfileManagerSignIn(username string, password string) (models.Profile, string, error) {
	var p models.Profile

	profileDb, message, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return p, message, err
	}

	if !models.ComparePassword(password, profileDb.PasswordEncrypt) {
		message := "bad password"
		return p, message, errors.New(message)
	}

	return profileDb, message, nil
}

func ProfileManagerSignUp(username string, password string) (models.Profile, string, error) {
	dateCreated := time.Now()
	p := models.Profile{}
	p.ProfileKey = p.NewProfileKey()
	p.Username = username
	hash, err := models.HashPassword(password)
	if err != nil {
		return models.Profile{}, "error while hashing password", err
	}

	p.PasswordEncrypt = hash
	p.DateCreated = dateCreated

	return repositories.ProfileCreate(p)
}

func ProfileManagerChangePassword(username string, password string, newPassword string) (models.Profile, string, error) {
	p, message, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return p, message, err
	}

	// check auth
	if !models.ComparePassword(password, p.PasswordEncrypt) {
		message := "bad password"
		return models.Profile{}, message, errors.New(message)
	}

	// check new password with current
	if models.ComparePassword(newPassword, p.PasswordEncrypt) {
		msg := "new password identical to the previous"
		return models.Profile{}, msg, errors.New(msg)
	}

	newPasswordEncrypt, err := models.HashPassword(newPassword)
	if err != nil {
		return models.Profile{}, "error while hashing new password", err
	}

	p.PasswordEncrypt = newPasswordEncrypt

	return repositories.ProfileUpdate(p)
}

func ProfileManagerDelete(username string, password string) (models.Profile, string, error) {
	p, message, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return p, message, err
	}

	// check auth
	if !models.ComparePassword(password, p.PasswordEncrypt) {
		message := "bad password"
		return models.Profile{}, message, errors.New(message)
	}

	return repositories.ProfileDelete(p)
}
