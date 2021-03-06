package managers

import (
	"errors"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
)

func ProfileManagerGet(profileKey string) (models.ProfileJson, string, error) {
	var pJson models.ProfileJson
	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return pJson, msg, err
	}

	return models.NewProfileJson(p), msg, nil
}

func ProfileManagerSignIn(username string, password string) (models.ProfileJson, string, error) {
	var p models.ProfileJson

	profileDb, message, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return p, message, err
	}

	if !models.ComparePassword(password, profileDb.PasswordEncrypt) {
		message := "bad password"
		return p, message, errors.New(message)
	}

	return models.NewProfileJson(profileDb), message, nil
}

func ProfileManagerChangePassword(username string, password string, newPassword string) (models.ProfileJson, string, error) {
	var pJson models.ProfileJson

	p, message, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return pJson, message, err
	}

	// check auth
	if !models.ComparePassword(password, p.PasswordEncrypt) {
		message := "bad password"
		return pJson, message, errors.New(message)
	}

	// check new password with current
	if models.ComparePassword(newPassword, p.PasswordEncrypt) {
		msg := "new password identical to the previous"
		return pJson, msg, errors.New(msg)
	}

	newPasswordEncrypt, err := models.HashPassword(newPassword)
	if err != nil {
		return pJson, "error while hashing new password", err
	}

	p.PasswordEncrypt = newPasswordEncrypt

	profileDb, msg, err := repositories.ProfileUpdate(p)
	if err != nil {
		logger.Error(err.Error())
		return pJson, msg, errors.New(msg)
	}

	return models.NewProfileJson(profileDb), "password updated", nil
}

func ProfileManagerDelete(username string, password string) (models.ProfileJson, string, error) {
	var pJson models.ProfileJson

	p, message, err := repositories.ProfileGetFromUsername(username)
	if err != nil {
		return pJson, message, err
	}

	// check auth
	if !models.ComparePassword(password, p.PasswordEncrypt) {
		message := "bad password"
		return pJson, message, errors.New(message)
	}

	profileDb, msg, err := repositories.ProfileDelete(p)
	if err != nil {
		logger.Error(err.Error())
		return pJson, msg, errors.New(msg)
	}

	return models.NewProfileJson(profileDb), "profile deleted", nil
}
