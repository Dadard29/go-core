package managers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Dadard29/go-core/api"
	"github.com/Dadard29/go-core/config"
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

func getQuota(profile models.Profile) (int, error) {
	var quotaStr string
	var err error
	if profile.Silver {
		quotaStr, err = api.Api.Config.GetValueFromFile(
			config.Profile,
			config.ProfileQuota,
			config.ProfileQuotaSilver)

		if err != nil {
			return 0, err
		}
	} else {
		quotaStr, err = api.Api.Config.GetValueFromFile(
			config.Profile,
			config.ProfileQuota,
			config.ProfileQuotaNotSilver)

		if err != nil {
			return 0, err
		}
	}

	quota, err := strconv.Atoi(quotaStr)
	if err != nil {
		return 0, err
	}

	return quota, nil
}

func SubsManagerGetFromToken(subToken string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	subDb, msg, err := repositories.SubsGetFromPkAndToken(subToken)
	if err != nil {
		return s, msg, err
	}

	profile, msg, err := repositories.ProfileGetFromKey(subDb.ProfileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(profile)
	if err != nil {
		return s, "error getting quota", err
	}

	if subDb.RequestCount >= quota {
		msg := fmt.Sprintf("quota (%d) reached", quota)
		return s, msg, errors.New(msg)
	}

	newSubDb, msg, err := repositories.SubsUpdateRequestCount(&subDb)
	if err != nil {
		return s, msg, err
	}

	a, msg, err := repositories.ApiGet(newSubDb.ApiName)
	if err != nil {
		return s, msg, err
	}

	return models.SubscriptionJson{
		AccessToken:    newSubDb.AccessToken,
		Api:            a,
		DateSubscribed: newSubDb.DateSubscribed,
		RequestCount:   newSubDb.RequestCount,
		Quota:          quota,
	}, "sub checked", nil
}

func SubsManagerGetFromApiName(apiName string, profileKey string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	profile, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(profile)
	if err != nil {
		return s, "error getting quota", err
	}

	subDb, msg, err := repositories.SubsGetFromApiName(apiName, profileKey)
	if err != nil {
		return s, msg, err
	}

	return models.SubscriptionJson{
		AccessToken:    subDb.AccessToken,
		Api:            a,
		DateSubscribed: subDb.DateSubscribed,
		RequestCount:   subDb.RequestCount,
		Quota:          quota,
	}, "sub checked", nil
}

func SubsManagerList(profileKey string) ([]models.SubscriptionJson, string, error) {
	var s []models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
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
			AccessToken:    sub.AccessToken,
			Api:            a,
			DateSubscribed: sub.DateSubscribed,
			RequestCount:   sub.RequestCount,
			Quota:          quota,
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

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	if p.Silver != true && a.Restricted {
		msg := "you need to be silver labelled to access restricted APIs"
		return s, msg, errors.New(msg)
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
		AccessToken:    subDb.AccessToken,
		Api:            a,
		DateSubscribed: subDb.DateSubscribed,
		RequestCount:   subDb.RequestCount,
		Quota:          quota,
	}

	return subJson, msg, nil
}

func SubsManagerDelete(profileKey string, apiName string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}

	subDb, msg, err := repositories.SubsDelete(p.ProfileKey, a.Name)
	if err != nil {
		return s, msg, err
	}

	subJson := models.SubscriptionJson{
		AccessToken:    subDb.AccessToken,
		Api:            a,
		DateSubscribed: subDb.DateSubscribed,
		RequestCount:   subDb.RequestCount,
		Quota:          quota,
	}

	return subJson, msg, nil
}

func SubsManagerUpdate(profileKey string, apiName string) (models.SubscriptionJson, string, error) {
	var s models.SubscriptionJson

	p, msg, err := repositories.ProfileGetFromKey(profileKey)
	if err != nil {
		return s, msg, err
	}

	quota, err := getQuota(p)
	if err != nil {
		return s, "error getting quota", err
	}

	a, msg, err := repositories.ApiGet(apiName)
	if err != nil {
		return s, msg, err
	}


	subUpdated, msg, err := repositories.SubsRegenerateToken(profileKey, apiName, subsGenerateAccessToken())
	if err != nil {
		return s, msg, err
	}

	subJson := models.SubscriptionJson{
		AccessToken:    subUpdated.AccessToken,
		Api:            a,
		DateSubscribed: subUpdated.DateSubscribed,
		RequestCount:   subUpdated.RequestCount,
		Quota:          quota,
	}

	return subJson, "token regenerated", nil
}
