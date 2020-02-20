package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/models"
)

func ProfileExists(p models.Profile) bool {
	var profileDb models.Profile
	api.Api.Database.Orm.Where(&models.Profile{Username:p.Username}).Find(&profileDb)
	return profileDb.Username == p.Username
}

func ProfileGetFromUsername(username string) (models.Profile, string, error) {
	var p models.Profile
	api.Api.Database.Orm.Where(&models.Profile{Username: username}).First(&p)

	if ! ProfileExists(p) {
		msg := fmt.Sprintf("no profile with username %s found", username)
		return p, msg, errors.New(msg)
	}

	return p, "profile retrieved", nil
}

func ProfileCreate(p models.Profile) (models.Profile, string, error) {
	if ProfileExists(p) {
		message := "existing profile with same username"
		return models.Profile{}, message, errors.New(message)
	}

	api.Api.Database.Orm.Create(&p)

	if ! ProfileExists(p) {
		message := "error creating profile"
		return models.Profile{}, message, errors.New(message)
	}

	return p, "profile created", nil
}

func ProfileUpdate(p models.Profile) (models.Profile, string, error) {
	if ! ProfileExists(p) {
		msg := fmt.Sprintf("no profile with username %s found", p.Username)
		return models.Profile{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Save(&p)

	return p, "profile updated", nil
}

func ProfileDelete(p models.Profile) (models.Profile, string, error) {
	if ! ProfileExists(p) {
		msg := fmt.Sprintf("no profile with username %s found", p.Username)
		return models.Profile{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Delete(&p)

	return p, "profile deleted", nil
}
