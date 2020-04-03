package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/models"
)

func profileExists(p models.Profile) bool {
	var profileDb models.Profile
	api.Api.Database.Orm.Where(&models.Profile{Username: p.Username}).Find(&profileDb)
	return profileDb.Username == p.Username
}

func ProfileGetFromKey(profileKey string) (models.Profile, string, error) {
	var p models.Profile
	api.Api.Database.Orm.Where(&models.Profile{
		ProfileKey: profileKey,
	}).First(&p)

	if !profileExists(p) {
		msg := fmt.Sprintf("no profile with this key found")
		return p, msg, errors.New(msg)
	}

	return p, "profile retrieved", nil
}

func ProfileGetFromUsername(username string) (models.Profile, string, error) {
	var p models.Profile
	api.Api.Database.Orm.Where(&models.Profile{Username: username}).First(&p)

	if !profileExists(p) {
		msg := fmt.Sprintf("no profile with username %s found", username)
		return p, msg, errors.New(msg)
	}

	return p, "profile retrieved", nil
}

func ProfileUpdate(p models.Profile) (models.Profile, string, error) {
	if !profileExists(p) {
		msg := fmt.Sprintf("no profile with username %s found", p.Username)
		return models.Profile{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Save(&p)

	return p, "profile updated", nil
}

func ProfileDelete(p models.Profile) (models.Profile, string, error) {
	if !profileExists(p) {
		msg := fmt.Sprintf("no profile with username %s found", p.Username)
		return models.Profile{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Delete(&p)

	return p, "profile deleted", nil
}
