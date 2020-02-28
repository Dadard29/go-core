package managers

import (
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"time"
)

func SubsManagerList(profileKey string) ([]models.SubscriptionJson, string, error) {
	var s []models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	subDb, msg, err := repositories.SubsListFromProfile(p)
	if err != nil {
		return s, msg, err
	}

	var subListJson = make([]models.SubscriptionJson, 0)
	for _, sub := range subDb {
		a, _, err := repositories.ApiGet(sub.ApiName)
		if err != nil {
			continue
		}

		subListJson = append(subListJson, models.SubscriptionJson{
			Profile:        models.NewProfileJson(p),
			Api:            a,
			DateSubscribed: sub.DateSubscribed,
		})
	}

	return subListJson, msg, nil
}

func SubsManagerCreate(profileKey string, apiName string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}
	
	subDb, msg, err := repositories.SubsCreate(models.Subscription{
		ProfileKey:     p.ProfileKey,
		ApiName:        a.Name,
		DateSubscribed: time.Now(),
	})
	if err != nil {
		return s, msg, err
	}

	subJson := models.SubscriptionJson{
		Profile:        models.NewProfileJson(p),
		Api:            a,
		DateSubscribed: subDb.DateSubscribed,
	}

	return subJson, msg, nil
}

func SubsManagerDelete(profileKey string, apiName string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	subDb, msg, err := repositories.SubsDelete(models.Subscription{
		ProfileKey: p.ProfileKey,
		ApiName:    a.Name,
	})
	if err != nil {
		return s, msg, err
	}

	subJson := models.SubscriptionJson{
		Profile:        models.NewProfileJson(p),
		Api:            a,
		DateSubscribed: subDb.DateSubscribed,
	}

	return subJson, msg, nil
}
