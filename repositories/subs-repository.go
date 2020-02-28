package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/models"
)

func subExists(s models.Subscription) bool {
	var subDb models.Subscription
	api.Api.Database.Orm.Where(&models.Subscription{
		ProfileKey: s.ProfileKey,
		ApiName:    s.ApiName,
	}).Find(&subDb)
	return subDb.ProfileKey == s.ProfileKey && subDb.ApiName == s.ApiName
}

func SubsExistsFromToken(subToken string) bool {
	var subDb models.Subscription
	api.Api.Database.Orm.Where(&models.Subscription{
		AccessToken: subToken,
	}).First(&subDb)
	return subDb.AccessToken == subToken
}

func SubsListFromProfile(p models.Profile) ([]models.Subscription, string, error) {
	var s []models.Subscription
	api.Api.Database.Orm.Find(&s, &models.Subscription{
		ProfileKey: p.ProfileKey,
	})

	return s, "subs listed", nil
}

// requires the fields sub.ProfileKey and sub.ApiName
func SubsCreate(s models.Subscription) (models.Subscription, string, error) {
	if subExists(s) {
		msg := fmt.Sprintf("existing sub for this user with api %s", s.ApiName)
		return models.Subscription{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Create(&s)

	if !subExists(s) {
		msg := "error subscribing"
		return models.Subscription{}, msg, errors.New(msg)
	}

	return s, "sub created", nil
}

// requires the fields sub.ProfileKey and sub.ApiName
func SubsDelete(s models.Subscription) (models.Subscription, string, error) {
	if !subExists(s) {
		msg := fmt.Sprintf("no sub found for this user with api %s", s.ApiName)
		return models.Subscription{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Delete(&s)

	if subExists(s) {
		msg := "error deleting sub"
		return models.Subscription{}, msg, errors.New(msg)
	}

	return s, "sub deleted", nil
}
