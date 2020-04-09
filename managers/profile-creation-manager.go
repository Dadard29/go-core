package managers

import (
	"errors"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"math/rand"
	"strconv"
	"time"
)

const (
	charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var confirmByMapping = map[string]func(contact string, tmpProfile models.TempProfile) error {
	"telegram": repositories.SendConfirmationTelegram,
	"email": repositories.SendConfirmationMail,
}

func generateConfirmationCode() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	code := ""
	for i := 0; i < 5; {
		randInt := rand.Intn(len(charset))
		c := string(charset[randInt])
		code += c
		i += 1
	}

	return code
}

func getExpirationDuration() time.Time {
	durationBeforeExpiration, _ := api.Api.Config.GetValueFromFile(
		config.Profile,
		config.ProfileCreation,
		config.ProfileCreationConfirmationExpirationDuration)

	durationBeforeExpirationInt, _ := strconv.Atoi(durationBeforeExpiration)

	expirationTime := time.Now().Add(
		time.Duration(durationBeforeExpirationInt) * time.Second)

	return expirationTime
}

func ConfirmationWayIsValid(confirmBy string) bool {
	_, c := confirmByMapping[confirmBy]
	return c
}

func ProfileManagerSendConfirmationCode(confirmBy string, contact string, tmpProfile models.TempProfile) (string, error) {
	if handler, c := confirmByMapping[confirmBy]; c {
		if err := handler(contact, tmpProfile); err != nil {
			return "error sending code", err
		}

		return "code send", nil

	} else {
		msg := "bad confirmation way"
		return msg, errors.New(msg)
	}
}

func ProfileManagerConfirmCode(username string, password string, confirmationCode string) (string, error) {
	tmpProfile, err := repositories.TempProfileGet(username)
	if err != nil {
		return "error getting temp profile", err
	}

	if tmpProfile.ExpirationTime.Before(time.Now()) {
		repositories.TempProfileDelete(username)
		msg := "code expired, ask for a new code"
		return msg, errors.New(msg)
	}

	if !models.ComparePassword(password, tmpProfile.PasswordEncrypt) {
		return "bad password given", err
	}

	if confirmationCode != tmpProfile.ConfirmationCode {
		msg := "bad code"
		return msg, errors.New(msg)
	}

	// code valid, account creation on going, to temp profile is useless
	repositories.TempProfileDelete(username)
	return "code checked", nil
}

func ProfileManagerCreateTemp(username string, password string) (models.TempProfile, string, error) {
	_, _, err := repositories.ProfileGetFromUsername(username)
	if err == nil {
		msg := "profile already existing with same username"
		return models.TempProfile{}, msg, errors.New(msg)
	}

	checkP, err := repositories.TempProfileGet(username)
	if err == nil {
		// not expired
		if checkP.ExpirationTime.After(time.Now()) {
			msg := "temporary profile already created"
			return models.TempProfile{}, msg, errors.New(msg)
		}

		repositories.TempProfileDelete(username)
		logger.Debug("deleting expired temp profile")
	}

	code := generateConfirmationCode()

	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		return models.TempProfile{}, "error hashing password", err
	}

	expirationTime := getExpirationDuration()

	p := models.TempProfile{
		Username:         username,
		PasswordEncrypt:  hashedPassword,
		ConfirmationCode: code,
		ExpirationTime:   expirationTime,
	}
	err = repositories.TempProfileCreate(p)
	if err != nil {
		return models.TempProfile{}, "error created temporary profile", err
	}

	return p, "temp profile created", nil
}

func ProfileManagerDeleteTemp(username string) error {
	_, err := repositories.TempProfileDelete(username)
	return err
}

func ProfileManagerCreate(username string, password string) (models.ProfileJson, string, error) {

	dateCreated := time.Now()
	p := models.Profile{}
	p.ProfileKey = p.NewProfileKey()
	p.Username = username
	hash, err := models.HashPassword(password)
	if err != nil {
		return models.ProfileJson{}, "error while hashing password", err
	}

	p.PasswordEncrypt = hash
	p.DateCreated = dateCreated
	p.Silver = false

	p.RecoverBy = ""
	p.Contact = ""

	profileDb, msg, err := repositories.ProfileCreate(p)
	if err != nil {
		logger.Error(err.Error())
		return models.ProfileJson{}, msg, errors.New(msg)
	}

	return models.NewProfileJson(profileDb), "profile created", nil
}
