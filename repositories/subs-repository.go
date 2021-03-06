package repositories

import (
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/models"
)

func subExistsFromPkAndApiName(s models.Subscription) bool {
	var subDb models.Subscription
	api.Api.Database.Orm.Where(&models.Subscription{
		ProfileKey:   s.ProfileKey,
		ApiName:      s.ApiName,
		FromEchoSlam: s.FromEchoSlam,
	}).Find(&subDb)
	return subDb.ProfileKey == s.ProfileKey && subDb.ApiName == s.ApiName && subDb.FromEchoSlam == s.FromEchoSlam
}

func SubsGetFromPkAndToken(subToken string) (models.Subscription, string, error) {
	var subDb models.Subscription
	api.Api.Database.Orm.Where(&models.Subscription{
		AccessToken: subToken,
	}).First(&subDb)
	if subDb.AccessToken == subToken {
		return subDb, "sub checked", nil
	} else {
		msg := "no sub for this user with this token"
		return models.Subscription{}, msg, errors.New(msg)
	}
}

func SubsGetFromApiName(apiName string, profileKey string, fromEchoSlam bool) (models.Subscription, string, error) {
	var subDb models.Subscription
	api.Api.Database.Orm.Where(&models.Subscription{
		ProfileKey:   profileKey,
		FromEchoSlam: fromEchoSlam,
		ApiName:      apiName,
	}).Find(&subDb)
	if subDb.ApiName == apiName && subDb.ProfileKey == profileKey && subDb.FromEchoSlam == fromEchoSlam {
		return subDb, "sub checked", nil
	} else {
		msg := "no sub with for this user and the api " + apiName
		return models.Subscription{}, msg, errors.New(msg)
	}
}

func SubsListFromProfile(p models.Profile) ([]models.Subscription, string, error) {
	var s []models.Subscription
	api.Api.Database.Orm.Find(&s, &models.Subscription{
		ProfileKey: p.ProfileKey,
		FromEchoSlam: false,
	})

	return s, "subs listed", nil
}

// requires the fields sub.ProfileKey and sub.ApiName
func SubsCreate(s models.Subscription) (models.Subscription, string, error) {
	if subExistsFromPkAndApiName(s) {
		msg := fmt.Sprintf("existing sub for this user with api %s", s.ApiName)
		return models.Subscription{}, msg, errors.New(msg)
	}

	api.Api.Database.Orm.Create(&s)

	if !subExistsFromPkAndApiName(s) {
		msg := "error subscribing"
		return models.Subscription{}, msg, errors.New(msg)
	}

	return s, "sub created", nil
}

// requires the fields sub.ProfileKey and sub.ApiName
func SubsDelete(profileKey string, apiName string) (models.Subscription, string, error) {
	if !subExistsFromPkAndApiName(models.Subscription{
		ProfileKey: profileKey,
		ApiName:    apiName,
	}) {
		msg := fmt.Sprintf("no sub found for this user with api %s", apiName)
		return models.Subscription{}, msg, errors.New(msg)
	}

	sDelete, msg, err := SubsGetFromApiName(apiName, profileKey, false)
	if err != nil {
		return models.Subscription{}, msg, err
	}

	api.Api.Database.Orm.Delete(&sDelete)

	if subExistsFromPkAndApiName(sDelete) {
		msg := "error deleting sub"
		return models.Subscription{}, msg, errors.New(msg)
	}

	return sDelete, "sub deleted", nil
}

func SubsRegenerateToken(profileKey string, apiName string, newAccessToken string) (models.Subscription, string, error) {
	if !subExistsFromPkAndApiName(models.Subscription{
		ProfileKey: profileKey,
		ApiName:    apiName,
	}) {
		msg := fmt.Sprintf("no sub found for this user with api %s", apiName)
		return models.Subscription{}, msg, errors.New(msg)
	}

	sUpdated, msg, err := SubsGetFromApiName(apiName, profileKey, false)
	if err != nil {
		return models.Subscription{}, msg, err
	}

	// deleting this sup
	api.Api.Database.Orm.Delete(&sUpdated)

	sUpdated.AccessToken = newAccessToken

	// creating a new one with updated token
	api.Api.Database.Orm.Create(&models.Subscription{
		AccessToken:    newAccessToken,
		ProfileKey:     sUpdated.ProfileKey,
		ApiName:        apiName,
		DateSubscribed: sUpdated.DateSubscribed,
		RequestCount:   sUpdated.RequestCount,
	})

	return sUpdated, "token regenerated", nil
}

func SubsUpdateRequestCount(subscription *models.Subscription) (*models.Subscription, string, error) {
	subscription.RequestCount += 1
	api.Api.Database.Orm.Save(subscription)

	return subscription, "request count incremented", nil
}

func SubsResetRequestCount(sub *models.Subscription) (*models.Subscription, string, error) {
	sub.RequestCount = 0
	api.Api.Database.Orm.Save(sub)

	return sub, "request count reset", nil
}

func SubsGetAll() []models.Subscription {
	var all []models.Subscription
	api.Api.Database.Orm.Find(&all)
	return all
}
