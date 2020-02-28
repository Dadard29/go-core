package managers

import (
	"crypto/sha256"
	"fmt"
	"github.com/Dadard29/go-core/models"
	"github.com/Dadard29/go-core/repositories"
	"math/rand"
	"strconv"
	"time"
)

func subsGenerateAccessToken() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	i := rand.Int()

	keyByte := []byte(strconv.Itoa(i))

	hash := sha256.New()

	hash.Write(keyByte)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func SubsManagerExists(subToken string) (bool, string, error) {
	return repositories.SubsExistsFromToken(subToken), "sub checked", nil
}

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
			AccessToken: 	sub.AccessToken,
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
		AccessToken:    subsGenerateAccessToken(),
		ProfileKey:     p.ProfileKey,
		ApiName:        a.Name,
		DateSubscribed: time.Now(),
	})
	if err != nil {
		return s, msg, err
	}

	subJson := models.SubscriptionJson{
		AccessToken: 	subDb.AccessToken,
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
		AccessToken: 	subDb.AccessToken,
		Profile:        models.NewProfileJson(p),
		Api:            a,
		DateSubscribed: subDb.DateSubscribed,
	}

	return subJson, msg, nil
}
